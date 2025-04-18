package services

import (
	"context"
	"fmt"
	"nexzap/internal/services/container"

	"github.com/docker/docker/client"
)

var tutorialsContainer = []container.Tutorial{
	{
		Name:      "0_go",
		Language:  "go",
		Image:     "gotest",
		WarmupDir: "",
		Command:   []string{"go", "test"},
	},
}

// Service encapsulates the state and operations for language testing services.
type Service struct {
	pool        container.Pool
	ctx         context.Context
	cli         *client.Client
	initialized bool
}

// NewService creates and initializes a new Service instance.
func NewService() (*Service, error) {
	svc := &Service{}
	if err := svc.init(); err != nil {
		return nil, err
	}
	return svc, nil
}

// init initializes the service, setting up the container pool and Docker client.
func (s *Service) init() error {
	s.pool = container.NewPool()
	s.ctx = context.Background()
	var err error
	s.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	s.initialized = true
	return nil
}

// RunTest executes the provided files in test mode for a given language.
func (s *Service) RunTest(number int, files []string) (string, error) {
	if !s.initialized {
		return "", fmt.Errorf("not initialized")
	}
	if len(tutorialsContainer) <= number {
		return "", fmt.Errorf("Number longer than number of tutorials")
	}
	language := tutorialsContainer[number]
	languagePool := s.pool.GetLanguagePool(s.ctx, s.cli, language)
	ctn, err := languagePool.GetContainer(s.ctx, s.cli)
	if err != nil {
		return "", err
	}
	output, err := container.Run(s.ctx, s.cli, ctn, files)
	languagePool.FreeContainer(s.ctx, s.cli, ctn)
	return output, err
}

// Cleanup stops and removes all containers in the pool.
func (s *Service) Cleanup() error {
	if !s.initialized {
		return fmt.Errorf("not initialized")
	}

	// Iterate through all language pools and clean up containers
	for _, tutorial := range tutorialsContainer {
		languagePool := s.pool.GetLanguagePool(s.ctx, s.cli, tutorial)
		// Drain minPool
		for len(languagePool.MinPool) > 0 {
			select {
			case ctn := <-languagePool.MinPool:
				container.StopAndRemove(s.ctx, s.cli, ctn)
			default:
				// No more containers to process
			}
		}
		// Drain extendedPool
		for len(languagePool.ExtendedPool) > 0 {
			select {
			case ctn := <-languagePool.ExtendedPool:
				container.StopAndRemove(s.ctx, s.cli, ctn)
			default:
				// No more containers to process
			}
		}
	}
	return nil
}
