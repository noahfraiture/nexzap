package container

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
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

	startTime := time.Now()
	log.Printf("Starting CopyToContainer for container %s at %v", ctn, startTime)
	err = cli.CopyToContainer(ctx, ctn, "/", archive, container.CopyToContainerOptions{})
	if err != nil {
		log.Printf("CopyToContainer for container %s failed, error: %v, duration: %v", ctn, err, time.Since(startTime))
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}
	log.Printf("CopyToContainer for container %s succeeded, duration: %v", ctn, time.Since(startTime))

	startTime = time.Now()
	log.Printf("Starting ContainerStart for container %s at %v", ctn, startTime)
	if err := cli.ContainerStart(ctx, ctn, container.StartOptions{}); err != nil {
		log.Printf("ContainerStart for container %s failed, error: %v, duration: %v", ctn, err, time.Since(startTime))
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}
	log.Printf("ContainerStart for container %s succeeded, duration: %v", ctn, time.Since(startTime))

	statusCh, errCh := cli.ContainerWait(ctx, ctn, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("ContainerWait for container %s failed with error: %v", ctn, err)
			stopAndRemove(ctx, cli, ctn)
			return "", err
		}
	case <-statusCh:
		log.Printf("ContainerWait for container %s completed", ctn)
	}

	logs, err := cli.ContainerLogs(ctx, ctn, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      startTime.Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("ContainerLogs for container %s failed, error: %v", ctn, err)
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}
	defer logs.Close()

	var logBytes bytes.Buffer
	_, err = io.Copy(&logBytes, logs)
	if err != nil {
		log.Printf("Copying logs for container %s failed, error: %v", ctn, err)
		stopAndRemove(ctx, cli, ctn)
		return "", err
	}

	return logBytes.String(), nil
}

// createContainer Creates a new container. If the WarmupDir is not empty, runs a warmup phase.
func createContainer(
	ctx context.Context,
	cli *client.Client,
	lang Language,
) (container.CreateResponse, error) {
	var emptyResp container.CreateResponse

	// Create container
	startTime := time.Now()
	log.Printf("Starting ContainerCreate for language %s at %v", lang.Name, startTime)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      lang.Image,
		Cmd:        lang.Command,
		WorkingDir: "/workspace",
		Tty:        false,
	}, nil, nil, nil, "")
	if err != nil {
		log.Printf("ContainerCreate for language %s failed, error: %v, duration: %v", lang.Name, err, time.Since(startTime))
		return emptyResp, err
	}
	log.Printf("ContainerCreate for language %s succeeded, container ID: %s, duration: %v", lang.Name, resp.ID, time.Since(startTime))

	// Warmup
	if lang.WarmupDir != "" {
		files, err := os.ReadDir(lang.WarmupDir)
		if err != nil {
			log.Printf("Reading WarmupDir %s for language %s failed, error: %v", lang.WarmupDir, lang.Name, err)
			stopAndRemove(ctx, cli, resp.ID)
			return emptyResp, err
		}

		filesName := make([]string, 0, len(files))
		for _, file := range files {
			filesName = append(filesName, lang.WarmupDir+file.Name())
		}

		log.Printf("Starting warmup phase for container %s (language: %s)", resp.ID, lang.Name)
		_, err = Run(ctx, cli, resp.ID, filesName)
		if err != nil {
			log.Printf("Warmup phase for container %s (language: %s) failed, error: %v", resp.ID, lang.Name, err)
			stopAndRemove(ctx, cli, resp.ID)
			return emptyResp, err
		}
		log.Printf("Warmup phase for container %s (language: %s) completed successfully", resp.ID, lang.Name)
	}
	return resp, nil
}

// stopAndRemove stops and removes a container.
func stopAndRemove(ctx context.Context, cli *client.Client, id string) {
	// Stop the container with a timeout of 10 seconds
	timeout := 10
	log.Printf("Stopping container %s with timeout %d seconds", id, timeout)
	if err := cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
		// If the container is already stopped or not found, proceed to remove it
		if !client.IsErrNotFound(err) {
			log.Printf("Failed to stop container %s, error: %v (proceeding to remove)", id, err)
			// Don't panic, just log the error and attempt removal
		} else {
			log.Printf("Container %s already stopped or not found, proceeding to remove", id)
		}
	} else {
		log.Printf("Container %s stopped successfully", id)
	}

	// Remove the container
	log.Printf("Removing container %s", id)
	if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{}); err != nil {
		if !client.IsErrNotFound(err) {
			log.Printf("Failed to remove container %s, error: %v", id, err)
			// Don't panic, just log the error
		} else {
			log.Printf("Container %s already removed or not found", id)
		}
	} else {
		log.Printf("Container %s removed successfully", id)
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
