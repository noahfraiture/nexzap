package container

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Run executes a container with the provided files and returns the output.
func Run(
	ctx context.Context,
	cli *client.Client,
	ctn string,
	files []File,
) (string, error) {
	archive, err := createTarArchive(files)
	if err != nil {
		StopAndRemove(ctx, cli, ctn)
		return "", err
	}

	err = cli.CopyToContainer(ctx, ctn, "/", archive, container.CopyToContainerOptions{})
	if err != nil {
		StopAndRemove(ctx, cli, ctn)
		return "", err
	}

	startTime := time.Now()
	if err := cli.ContainerStart(ctx, ctn, container.StartOptions{}); err != nil {
		StopAndRemove(ctx, cli, ctn)
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ctn, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			StopAndRemove(ctx, cli, ctn)
			return "", err
		}
	case <-statusCh:
	}

	logs, err := cli.ContainerLogs(ctx, ctn, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      startTime.Format(time.RFC3339),
	})
	if err != nil {
		StopAndRemove(ctx, cli, ctn)
		return "", err
	}
	defer logs.Close()

	var logBytes bytes.Buffer
	_, err = io.Copy(&logBytes, logs)
	if err != nil {
		StopAndRemove(ctx, cli, ctn)
		return "", err
	}

	return logBytes.String(), nil
}

// createContainer Creates a new container. If the WarmupDir is not empty, runs a warmup phase.
func createContainer(
	ctx context.Context,
	cli *client.Client,
	lang Tutorial,
) (container.CreateResponse, error) {
	var emptyResp container.CreateResponse

	// Create container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      lang.Image,
		Cmd:        lang.Command,
		WorkingDir: "/workspace",
		Tty:        false,
	}, nil, nil, nil, "")
	if err != nil {
		return emptyResp, err
	}

	return resp, nil
}

// StopAndRemove stops and removes a container.
func StopAndRemove(ctx context.Context, cli *client.Client, id string) {
	// Stop the container with a timeout of 10 seconds
	timeout := 10
	if err := cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
		// If the container is already stopped or not found, proceed to remove it
		if !client.IsErrNotFound(err) {
			// TODO : change that
			panic(err)
			// Don't panic, just log the error and attempt removal
		}
	}

	// Remove the container
	if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{}); err != nil {
		if !client.IsErrNotFound(err) {
			panic(err)
			// Don't panic, just log the error
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
