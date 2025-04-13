package container

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/client"
)

const (
	// number of always-up container
	MIN_CTN = 3
	// number of maximum additional container
	MAX_CTN = 10
	// time before language is fully discarded and thus base container
	LANGUAGE_TIMEOUT = 3 * time.Minute
	// time before a extended container is discared
	CONTAINER_TIMEOUT = 15 * time.Second
	// time before we signal there's no container available
	WAIT_TIMEOUT = 10 * time.Second
	// number of container of margin
	CONTAINER_MARGIN = 2
)

type Tutorial struct {
	Name      string
	Language  string
	Image     string
	WarmupDir string
	Command   []string
}

type Pool struct {
	sync.Mutex
	pool map[string]LanguagePool
}

func NewPool() Pool {
	return Pool{pool: make(map[string]LanguagePool)}
}

// GetLanguagePool creates a pool for a given language if it doesn't exist, then returns it.
// Synchronized method to avoid duplicate language pool.
func (p *Pool) GetLanguagePool(ctx context.Context, cli *client.Client, language Tutorial) LanguagePool {
	p.Lock()
	defer p.Unlock()
	if lp, ok := p.pool[language.Language]; ok {
		return lp
	}
	p.newLanguage(ctx, cli, language)
	return p.pool[language.Language]
}

// newLanguage adds a pool for a language in the main pool.
// Defines its configuration and starts its base container.
func (p *Pool) newLanguage(
	ctx context.Context,
	cli *client.Client,
	lang Tutorial,
) {
	// Create the language
	language := LanguagePool{
		minPool:         make(chan string, MIN_CTN),
		extendedPool:    make(chan string, MAX_CTN),
		available:       make(chan any, MAX_CTN),
		language:        lang,
		languageTimeout: *NewTimeout(LANGUAGE_TIMEOUT, nil),
		extendTimeout:   *NewTimeout(CONTAINER_TIMEOUT, nil),
	}

	for range MAX_CTN {
		language.available <- struct{}{}
	}

	language.languageTimeout.action = func() {
		p.cleanLanguage(ctx, cli, lang.Language) // FIX: clean can try to remove self?
	}

	language.extendTimeout.action = func() {
		var wg sync.WaitGroup
		for c := range language.extendedPool {
			wg.Add(1)
			go func(containerID string) {
				defer wg.Done()
				stopAndRemove(ctx, cli, containerID)
			}(c)
		}
		wg.Wait()
		for range MAX_CTN - len(language.available) {
			language.available <- struct{}{}
		}
	}

	// Create the minPool
	var wg sync.WaitGroup
	wg.Add(MIN_CTN)
	for range MIN_CTN {
		go func() {
			defer wg.Done()
			resp, err := createContainer(ctx, cli, lang)
			if err != nil {
				panic(err)
			}
			language.minPool <- resp.ID
		}()
	}
	wg.Wait()

	p.pool[lang.Language] = language
}

// LanguagePool represents a pool of containers for a specific language.
type LanguagePool struct {
	language        Tutorial
	minPool         chan string // pool of containers that should always be running
	extendedPool    chan string // pool of containers that can shrink or expand
	available       chan any    // quantity of containers still possible to deploy
	languageTimeout Timeout
	extendTimeout   Timeout
}

// GetContainer queries a container from the language pool and resets the timeout.
// The language pool can create new containers to keep a margin.
// You must free it after usage.
func (lp *LanguagePool) GetContainer(ctx context.Context, cli *client.Client) (string, error) {
	select {
	case c := <-lp.minPool:
		lp.languageTimeout.StartTimer()
		extendContainer(ctx, cli, lp)
		return c, nil
	case c := <-lp.extendedPool:
		lp.languageTimeout.StartTimer()
		lp.extendTimeout.StartTimer()
		extendContainer(ctx, cli, lp)
		return c, nil
	case <-time.After(WAIT_TIMEOUT):
		return "", fmt.Errorf("timeout waiting for container")
	}
}

// FreeContainer returns a container to the pool after usage.
func (lp *LanguagePool) FreeContainer(
	ctx context.Context,
	cli *client.Client,
	ctn string,
) {
	select {
	case lp.minPool <- ctn:
	case lp.extendedPool <- ctn:
	default:
		stopAndRemove(ctx, cli, ctn)
	}
}

// extendContainer extends the container pool if necessary to maintain a margin.
func extendContainer(ctx context.Context, cli *client.Client, lp *LanguagePool) {
	nbFree := len(lp.extendedPool) + len(lp.minPool)
	if nbFree < CONTAINER_MARGIN {
		select {
		case <-lp.available:
			go createAndAddContainer(ctx, cli, lp)
		default:
		}
	}
}

// createAndAddContainer creates a new container and adds it to the extended pool.
func createAndAddContainer(ctx context.Context, cli *client.Client, lp *LanguagePool) {
	resp, err := createContainer(ctx, cli, lp.language)
	if err != nil {
		lp.available <- struct{}{}
		return
	}
	lp.extendedPool <- resp.ID
}

// cleanLanguage removes a language from the main pool and cleans the minPool.
// extendedPool will automatically be cleaned with a proper timeout.
func (p *Pool) cleanLanguage(ctx context.Context, cli *client.Client, name string) {
	// TODO: fix that?
	language, ok := p.pool[name]
	if !ok {
		return
	}
	close(language.minPool)
	close(language.extendedPool)
	for ctn := range language.minPool {
		stopAndRemove(ctx, cli, ctn)
	}
}
