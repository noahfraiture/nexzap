package container

import (
	"context"
	"fmt"
	"log"
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

type Language struct {
	Name      string
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
func (p *Pool) GetLanguagePool(ctx context.Context, cli *client.Client, language Language) LanguagePool {
	p.Lock()
	defer p.Unlock()
	if lp, ok := p.pool[language.Name]; ok {
		return lp
	}
	p.newLanguage(ctx, cli, language)
	return p.pool[language.Name]
}

// newLanguage adds a pool for a language in the main pool.
// Defines its configuration and starts its base container.
func (p *Pool) newLanguage(
	ctx context.Context,
	cli *client.Client,
	lang Language,
) {
	log.Printf("Creating new language pool for %s", lang.Name)
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
		p.cleanLanguage(ctx, cli, lang.Name) // FIX: clean can try to remove self?
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

	p.pool[lang.Name] = language
}

// LanguagePool represents a pool of containers for a specific language.
type LanguagePool struct {
	language        Language
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
	log.Printf("Attempting to get container for language %s", lp.language.Name)
	selectStart := time.Now()
	select {
	case c := <-lp.minPool:
		lp.languageTimeout.StartTimer()
		log.Printf("Container %s taken from minPool for language %s, select took %v", c[:12], lp.language.Name, time.Since(selectStart))
		extendContainer(ctx, cli, lp)
		return c, nil
	case c := <-lp.extendedPool:
		lp.languageTimeout.StartTimer()
		lp.extendTimeout.StartTimer()
		log.Printf("Container %s taken from extendedPool for language %s, select took %v", c[:12], lp.language.Name, time.Since(selectStart))
		extendContainer(ctx, cli, lp)
		return c, nil
	case <-time.After(WAIT_TIMEOUT):
		log.Printf("Timeout waiting for container for language %s, select took %v", lp.language.Name, time.Since(selectStart))
		return "", fmt.Errorf("timeout waiting for container")
	}
}

// FreeContainer returns a container to the pool after usage.
func (lp *LanguagePool) FreeContainer(
	ctx context.Context,
	cli *client.Client,
	ctn string,
) {
	log.Printf("Attempting to free container %s for language %s", ctn[:12], lp.language.Name)
	selectStart := time.Now()
	select {
	case lp.minPool <- ctn:
		log.Printf("Container %s returned to minPool, select took %v", ctn[:12], time.Since(selectStart))
	case lp.extendedPool <- ctn:
		log.Printf("Container %s returned to extendedPool, select took %v", ctn[:12], time.Since(selectStart))
	default:
		log.Printf("Container %s could not be returned to any pool, will be stopped and removed, select took %v", ctn[:12], time.Since(selectStart))
		stopAndRemove(ctx, cli, ctn)
	}
}

// extendContainer extends the container pool if necessary to maintain a margin.
func extendContainer(ctx context.Context, cli *client.Client, lp *LanguagePool) {
	nbFree := len(lp.extendedPool) + len(lp.minPool)
	if nbFree < CONTAINER_MARGIN {
		log.Printf("Attempting to extend container pool for language %s", lp.language.Name)
		selectStart := time.Now()
		select {
		case <-lp.available:
			log.Printf("Acquired slot to extend container for language %s, select took %v", lp.language.Name, time.Since(selectStart))
			go createAndAddContainer(ctx, cli, lp)
		default:
			log.Printf("No available slots to extend container for language %s, select took %v", lp.language.Name, time.Since(selectStart))
		}
	}
}

// createAndAddContainer creates a new container and adds it to the extended pool.
func createAndAddContainer(ctx context.Context, cli *client.Client, lp *LanguagePool) {
	resp, err := createContainer(ctx, cli, lp.language)
	if err != nil {
		lp.available <- struct{}{}
		log.Printf("Failed to extend container for language %s: %v", lp.language.Name, err)
		return
	}
	lp.extendedPool <- resp.ID
	log.Printf("Extended container pool with new container %s for language %s", resp.ID[:12], lp.language.Name)
}

// cleanLanguage removes a language from the main pool and cleans the minPool.
// extendedPool will automatically be cleaned with a proper timeout.
func (p *Pool) cleanLanguage(ctx context.Context, cli *client.Client, name string) {
	// TODO: fix that?
	language, ok := p.pool[name]
	if !ok {
		log.Printf("Attempted to clean non-existent language pool: %s", name)
		return
	}
	log.Printf("Cleaning language pool for %s", name)
	close(language.minPool)
	close(language.extendedPool)
	for ctn := range language.minPool {
		stopAndRemove(ctx, cli, ctn)
	}
}
