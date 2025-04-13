package container

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Run executes a container with the provided files and returns the output.
func Run(
	ctx context.Context,
	cli *client.Client,
	ctn string,
	files []string,
) (string, error) {
	archive, err := createTarArchive(files)
	if err != nil {
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}

	err = cli.CopyToContainer(ctx, ctn, "/", archive, container.CopyToContainerOptions{})
	if err != nil {
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}

	startTime := time.Now()
	if err := cli.ContainerStart(ctx, ctn, container.StartOptions{}); err != nil {
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ctn, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			stopAndRemove(ctx, cli, ctn)
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
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}
	defer logs.Close()

	var logBytes bytes.Buffer
	_, err = io.Copy(&logBytes, logs)
	if err != nil {
		stopAndRemove(ctx, cli, ctn)
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

	// Warmup
	if lang.WarmupDir != "" {
		files, err := os.ReadDir(lang.WarmupDir)
		if err != nil {
			stopAndRemove(ctx, cli, resp.ID)
			return emptyResp, err
		}

		filesName := make([]string, 0, len(files))
		for _, file := range files {
			filesName = append(filesName, lang.WarmupDir+file.Name())
		}

		_, err = Run(ctx, cli, resp.ID, filesName)
		if err != nil {
			stopAndRemove(ctx, cli, resp.ID)
			return emptyResp, err
		}
	}
	return resp, nil
}

// stopAndRemove stops and removes a container.
func stopAndRemove(ctx context.Context, cli *client.Client, id string) {
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

// createTarArchive creates a tar archive from the provided files.
func createTarArchive(files []string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		header := &tar.Header{
			Name:    "workspace/" + filepath.Base(filePath),
			Size:    stat.Size(),
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
		}

		if err := tw.WriteHeader(header); err != nil {
			return nil, err
		}

		if _, err := io.Copy(tw, file); err != nil {
			return nil, err
		}
	}

	return &buf, nil
}
