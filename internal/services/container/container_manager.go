package container

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type RunResponse = container.WaitResponse

// Run executes a container with the provided files and returns the output.
func Run(
	ctx context.Context,
	cli *client.Client,
	ctn string,
	files []File,
) (string, RunResponse, error) {
	var err error
	defer func() {
		if err != nil {
			StopAndRemove(ctx, cli, ctn)
		}
	}()
	archive, err := createTarArchive(files)
	if err != nil {
		return "", RunResponse{}, err
	}

	err = cli.CopyToContainer(ctx, ctn, "/", archive, container.CopyToContainerOptions{})
	if err != nil {
		return "", RunResponse{}, err
	}

	startTime := time.Now()
	if err := cli.ContainerStart(ctx, ctn, container.StartOptions{}); err != nil {
		return "", RunResponse{}, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ctn, container.WaitConditionNotRunning)
	var status container.WaitResponse
	select {
	case err := <-errCh:
		if err != nil {
			return "", RunResponse{}, err
		}
	case status = <-statusCh:
	}

	logs, err := cli.ContainerLogs(ctx, ctn, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      startTime.Format(time.RFC3339),
	})
	if err != nil {
		return "", RunResponse{}, err
	}
	defer logs.Close()

	var logBytes bytes.Buffer
	_, err = io.Copy(&logBytes, logs)
	if err != nil {
		return "", RunResponse{}, err
	}

	return logBytes.String(), status, nil
}

// createContainer Creates a new container. If the WarmupDir is not empty, runs a warmup phase.
func createContainer(
	ctx context.Context,
	cli *client.Client,
	lang Tutorial,
) (container.CreateResponse, error) {
	var emptyResp container.CreateResponse

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      lang.Image,
		Cmd:        lang.Command,
		WorkingDir: "/workspace",
		Tty:        false,
	}, &container.HostConfig{
		CapDrop: []string{"ALL"},
		CapAdd:  []string{},
		// Prevent mounting the Docker socket or other sensitive paths
		Binds: nil,
		Resources: container.Resources{
			Memory:    512 * 1024 * 1024,
			CPUQuota:  100000,
			CPUPeriod: 100000,
		},
		// Isolate from host network for additional security
		NetworkMode: "none",
	}, nil, nil, "")
	if err != nil {
		fmt.Println(err)
		StopAndRemove(ctx, cli, resp.ID) // in case the error still creates the container
		return emptyResp, err
	}

	return resp, nil
}

// StopAndRemove stops and removes a container.
func StopAndRemove(ctx context.Context, cli *client.Client, id string) {
	timeout := 10
	if err := cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
		if !client.IsErrNotFound(err) {
			log.Println(err)
		}
	}

	if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{}); err != nil {
		if !client.IsErrNotFound(err) {
			log.Println(err)
		}
	}
}

type File struct {
	Name    string
	Content string
}

// createTarArchive creates a tar archive from the provided files.
func createTarArchive(files []File) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	for _, file := range files {
		header := &tar.Header{
			Name:    "workspace/" + file.Name,
			Size:    int64(len(file.Content)),
			Mode:    0644,
			ModTime: time.Now(),
		}

		if err := tw.WriteHeader(header); err != nil {
			return nil, err
		}

		if _, err := io.WriteString(tw, file.Content); err != nil {
			return nil, err
		}
	}

	return &buf, nil
}
