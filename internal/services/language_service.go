package services

import (
	"context"
	"fmt"
	"log"
	"zapbyte/internal/services/container"

	"github.com/docker/docker/client"
)

type LanguageName int

const (
	_ LanguageName = iota
	GO
)

var languages = map[LanguageName]container.Language{
	GO: {
		Name:      "go",
		Image:     "gotest",
		WarmupDir: "languages/go/",
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
func (s *Service) RunTest(name LanguageName, files []string) (string, error) {
	if !s.initialized {
		return "", fmt.Errorf("not initialized")
	}
	if _, ok := languages[name]; !ok {
		return "", fmt.Errorf("unknown language %d", name)
	}
	language := languages[name]
	languagePool := s.pool.GetLanguagePool(s.ctx, s.cli, language)
	ctn, err := languagePool.GetContainer(s.ctx, s.cli)
	if err != nil {
		return "", err
	}
	log.Printf("Starting tests on container %s", ctn[:12])
	output, err := container.Run(s.ctx, s.cli, ctn, files)
	languagePool.FreeContainer(s.ctx, s.cli, ctn)
	return output, err
}
