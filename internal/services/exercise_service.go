package services

import (
	"context"
	"fmt"
	"nexzap/internal/services/container"
	"strings"

	generated "nexzap/internal/db/generated"

	"github.com/docker/docker/client"
)

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

type Correction = generated.FindSubmissionDataRow

// RunTest executes the provided files in test mode for a given language.
func (s *Service) RunTest(correction Correction, payload string) (string, error) {
	if !s.initialized {
		return "", fmt.Errorf("not initialized")
	}
	tutorial := container.Tutorial{
		Image:   correction.DockerImage,
		Command: strings.Split(correction.Command, " "),
	}
	languagePool := s.pool.GetLanguagePool(s.ctx, s.cli, tutorial)
	ctn, err := languagePool.GetContainer(s.ctx, s.cli)
	if err != nil {
		return "", err
	}

	files := []container.File{}
	for i := range correction.FilesName {
		files = append(files, container.File{
			Name:    correction.FilesName[i],
			Content: correction.FilesContent[i],
		})
	}
	files = append(files, container.File{
		Name:    correction.SubmissionName,
		Content: payload,
	})
	output, err := container.Run(s.ctx, s.cli, ctn, files)
	languagePool.FreeContainer(s.ctx, s.cli, ctn)
	return output, err
}

// Cleanup stops and removes all containers in the pool.
func (s *Service) Cleanup() error {
	if !s.initialized {
		return fmt.Errorf("not initialized")
	}

	s.pool.CleanAll(s.ctx, s.cli)
	return nil
}
