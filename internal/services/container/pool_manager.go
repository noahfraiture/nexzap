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
	Image   string
	Command []string
}

type Pool struct {
	sync.Mutex
	pool map[string]ImagePool
}

func NewPool() Pool {
	return Pool{pool: make(map[string]ImagePool)}
}

// GetLanguagePool creates a pool for a given language if it doesn't exist, then returns it.
// Synchronized method to avoid duplicate language pool.
func (p *Pool) GetLanguagePool(ctx context.Context, cli *client.Client, language Tutorial) ImagePool {
	p.Lock()
	defer p.Unlock()
	if lp, ok := p.pool[language.Image]; ok {
		return lp
	}
	p.newImage(ctx, cli, language)
	return p.pool[language.Image]
}

// newImage adds a pool for a language in the main pool.
// Defines its configuration and starts its base container.
func (p *Pool) newImage(
	ctx context.Context,
	cli *client.Client,
	lang Tutorial,
) {
	// Create the language
	language := ImagePool{
		MinPool:         make(chan string, MIN_CTN),
		ExtendedPool:    make(chan string, MAX_CTN),
		available:       make(chan any, MAX_CTN),
		language:        lang,
		languageTimeout: *NewTimeout(LANGUAGE_TIMEOUT, nil),
		extendTimeout:   *NewTimeout(CONTAINER_TIMEOUT, nil),
	}

	for range MAX_CTN {
		language.available <- struct{}{}
	}

	language.languageTimeout.action = func() {
		p.cleanLanguage(ctx, cli, lang.Image) // FIX: clean can try to remove self?
	}

	language.extendTimeout.action = func() {
		var wg sync.WaitGroup
		for c := range language.ExtendedPool {
			wg.Add(1)
			go func(containerID string) {
				defer wg.Done()
				StopAndRemove(ctx, cli, containerID)
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
			language.MinPool <- resp.ID
		}()
	}
	wg.Wait()

	p.pool[lang.Image] = language
}

// ImagePool represents a pool of containers for a specific language.
type ImagePool struct {
	language        Tutorial
	MinPool         chan string // pool of containers that should always be running
	ExtendedPool    chan string // pool of containers that can shrink or expand
	available       chan any    // quantity of containers still possible to deploy
	languageTimeout Timeout
	extendTimeout   Timeout
}

// GetContainer queries a container from the language pool and resets the timeout.
// The language pool can create new containers to keep a margin.
// You must free it after usage.
func (lp *ImagePool) GetContainer(ctx context.Context, cli *client.Client) (string, error) {
	select {
	case c := <-lp.MinPool:
		lp.languageTimeout.StartTimer()
		extendContainer(ctx, cli, lp)
		return c, nil
	case c := <-lp.ExtendedPool:
		lp.languageTimeout.StartTimer()
		lp.extendTimeout.StartTimer()
		extendContainer(ctx, cli, lp)
		return c, nil
	case <-time.After(WAIT_TIMEOUT):
		return "", fmt.Errorf("timeout waiting for container")
	}
}

// FreeContainer returns a container to the pool after usage.
func (lp *ImagePool) FreeContainer(
	ctx context.Context,
	cli *client.Client,
	ctn string,
) {
	select {
	case lp.MinPool <- ctn:
	case lp.ExtendedPool <- ctn:
	default:
		StopAndRemove(ctx, cli, ctn)
	}
}

// extendContainer extends the container pool if necessary to maintain a margin.
func extendContainer(ctx context.Context, cli *client.Client, lp *ImagePool) {
	nbFree := len(lp.ExtendedPool) + len(lp.MinPool)
	if nbFree < CONTAINER_MARGIN {
		select {
		case <-lp.available:
			go createAndAddContainer(ctx, cli, lp)
		default:
		}
	}
}

// createAndAddContainer creates a new container and adds it to the extended pool.
func createAndAddContainer(ctx context.Context, cli *client.Client, lp *ImagePool) {
	resp, err := createContainer(ctx, cli, lp.language)
	if err != nil {
		lp.available <- struct{}{}
		return
	}
	lp.ExtendedPool <- resp.ID
}

// cleanLanguage removes a language from the main pool and cleans the minPool.
// extendedPool will automatically be cleaned with a proper timeout.
func (p *Pool) cleanLanguage(ctx context.Context, cli *client.Client, name string) {
	// TODO: fix that?
	language, ok := p.pool[name]
	if !ok {
		return
	}
	close(language.MinPool)
	close(language.ExtendedPool)
	for ctn := range language.MinPool {
		StopAndRemove(ctx, cli, ctn)
	}
}

func (p *Pool) CleanAll(ctx context.Context, cli *client.Client) {
	p.Lock()
	defer p.Unlock()

	for imageName := range p.pool {
		p.cleanLanguage(ctx, cli, imageName)
		delete(p.pool, imageName)
	}
}
