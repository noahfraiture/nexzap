package services

import (
	"context"
	"fmt"
	"nexzap/internal/services/container"
	"strings"
	"time"

	generated "nexzap/internal/db/generated"

	"github.com/docker/docker/client"
)

// ExerciseService encapsulates the state and operations for language testing services.
type ExerciseService struct {
	pool        container.Pool
	ctx         context.Context
	cli         *client.Client
	initialized bool
}

// NewExerciseService creates and initializes a new Service instance.
func NewExerciseService() (*ExerciseService, error) {
	svc := &ExerciseService{}
	if err := svc.init(); err != nil {
		return nil, err
	}
	return svc, nil
}

// init initializes the service, setting up the container pool and Docker client.
func (s *ExerciseService) init() error {
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

// TODO : message clietn side on error in running
// RunTest executes the provided files in test mode for a given language.
func (s *ExerciseService) RunTest(correction Correction, payload string) (string, container.RunResponse, error) {
	if !s.initialized {
		return "", container.RunResponse{}, fmt.Errorf("not initialized")
	}
	tutorial := container.Tutorial{
		Image:   correction.DockerImage,
		Command: strings.Split(correction.Command, " "),
	}

	files := []container.File{}
	for i, name := range correction.FilesName {
		if name != correction.SubmissionName {
			files = append(files, container.File{
				Name:    correction.FilesName[i],
				Content: correction.FilesContent[i],
			})
		}
	}
	files = append(files, container.File{
		Name:    correction.SubmissionName,
		Content: payload,
	})

	languagePool := s.pool.GetImagePool(s.ctx, s.cli, tutorial)
	ctn, err := languagePool.GetContainer(s.ctx, s.cli)
	if err != nil {
		return "", container.RunResponse{}, err
	}
	for attempt := range 3 {
		timeoutCtx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
		defer cancel()

		output, status, err := container.Run(timeoutCtx, s.cli, ctn, files)
		languagePool.FreeContainer(s.ctx, s.cli, ctn)
		if err == nil {
			return output, status, nil
		}
		fmt.Printf("Attempt %d failed: %v\n", attempt, err)
		ctn, err = languagePool.GetContainer(s.ctx, s.cli)
		if err != nil {
			return "", container.RunResponse{}, err
		}
		// Wait before retrying, but only if this isn't the last attempt
		if attempt < 3 {
			time.Sleep(3 * time.Second)
		}
	}
	return "", container.RunResponse{}, fmt.Errorf("unexpected error in retry loop")
}

// Cleanup stops and removes all containers in the pool.
func (s *ExerciseService) Cleanup() error {
	if !s.initialized {
		return fmt.Errorf("not initialized")
	}

	s.pool.CleanAll(s.ctx, s.cli)
	return nil
}
