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

// GetImagePool creates a pool for a given language if it doesn't exist, then returns it.
// Synchronized method to avoid duplicate language pool.
func (p *Pool) GetImagePool(ctx context.Context, cli *client.Client, tutorial Tutorial) ImagePool {
	p.Lock()
	defer p.Unlock()
	if lp, ok := p.pool[tutorial.Image]; ok {
		return lp
	}
	p.newImage(ctx, cli, tutorial)
	return p.pool[tutorial.Image]
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
		extensionSlots:  make(chan any, MAX_CTN),
		language:        lang,
		languageTimeout: *NewTimeout(LANGUAGE_TIMEOUT, nil),
		extendTimeout:   *NewTimeout(CONTAINER_TIMEOUT, nil),
	}

	language.languageTimeout.action = func() {
		p.cleanImage(ctx, cli, lang.Image)
	}

	language.extendTimeout.action = func() {
		for c := range language.ExtendedPool {
			StopAndRemove(ctx, cli, c)
			<-language.extensionSlots
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
	language Tutorial
	// pool of containers that should always be running
	MinPool chan string
	// pool of containers that can shrink or expand
	ExtendedPool chan string
	// Quantity of containers still possible to deploy.
	// Needed to take place before instanciating the container in ExtendedPool
	extensionSlots  chan any
	languageTimeout Timeout
	extendTimeout   Timeout
}

// GetContainer queries a container from the language pool and resets the timeout.
// The language pool can create new containers to keep a margin.
// You must free it after usage.
func (lp *ImagePool) GetContainer(ctx context.Context, cli *client.Client) (string, error) {
	// Check if container in MinPool
	select {
	case c := <-lp.MinPool:
		lp.languageTimeout.StartTimer()
		return c, nil
	default:
	}

	// Check if container in ExtendedPool
	select {
	case c := <-lp.ExtendedPool:
		lp.languageTimeout.StartTimer()
		lp.extendTimeout.StartTimer()
		return c, nil
	default:
	}

	// NOTE : the container hitting the limit of available container will always
	// have to wait for a new one. This slows down the system with burst of
	// submission. However we can consider this as acceptable for simplicity

	// Else has too wait and get first free container
	extendContainer(ctx, cli, lp)
	select {
	case c := <-lp.MinPool:
		lp.languageTimeout.StartTimer()
		return c, nil
	case c := <-lp.ExtendedPool:
		lp.languageTimeout.StartTimer()
		lp.extendTimeout.StartTimer()
		return c, nil
	case <-time.After(WAIT_TIMEOUT):
		return "", fmt.Errorf("timeout waiting for container")
	}

}

// FreeContainer returns a container to the pool after usage.
// WARNING : it does not stop or remove the container
func (lp *ImagePool) FreeContainer(
	ctx context.Context,
	cli *client.Client,
	ctn string,
) {
	// First try to give it to MinPool
	select {
	case lp.MinPool <- ctn:
		return
	default:
	}

	// Else give where there's space left
	select {
	case lp.MinPool <- ctn:
	case lp.ExtendedPool <- ctn:
	default:
		log.Fatalf("Container that shouldn't have been up")
	}
}

// extendContainer extends the container pool if there's still slot available
func extendContainer(ctx context.Context, cli *client.Client, lp *ImagePool) {
	select {
	case lp.extensionSlots <- struct{}{}:
		go createAndAddContainer(ctx, cli, lp)
	default:
	}
}

// createAndAddContainer creates a new container and adds it to the extended pool.
func createAndAddContainer(ctx context.Context, cli *client.Client, lp *ImagePool) {
	resp, err := createContainer(ctx, cli, lp.language)
	if err != nil {
		<-lp.extensionSlots
	} else {
		lp.ExtendedPool <- resp.ID
	}
}

// cleanImage removes a language from the main pool and cleans the minPool.
// extendedPool will automatically be cleaned with a proper timeout.
func (p *Pool) cleanImage(ctx context.Context, cli *client.Client, name string) {
	p.Lock()
	defer p.Unlock()
	language, ok := p.pool[name]
	if !ok {
		return
	}
	close(language.MinPool)
	close(language.ExtendedPool)
	for ctn := range language.MinPool {
		StopAndRemove(ctx, cli, ctn)
	}
	delete(p.pool, name)
}

// CleanAll is not concurrent safe !
func (p *Pool) CleanAll(ctx context.Context, cli *client.Client) {
	for imageName := range p.pool {
		p.cleanImage(ctx, cli, imageName)
		delete(p.pool, imageName)
	}
}
