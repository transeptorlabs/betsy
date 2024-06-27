package docker

import (
	"fmt"
	"os"
	"strings"

	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog/log"
)

// ContainerManager manages containers
type ContainerManager struct {
	supportedImages map[string]ContainerDetails
	client          *client.Client
	ctx             context.Context
}

// ContainerDetails contains details of a container
type ContainerDetails struct {
	imageName     string
	containerName string
	ContainerID   string
	IsRunning     bool
	Cmd           []string
	ExposedPorts  nat.PortSet
}

// NewContainerManagerr creates a new container manager
func NewContainerManager() (*ContainerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &ContainerManager{
		supportedImages: map[string]ContainerDetails{
			"transeptor": {
				containerName: "betsy-transeptor",
				ContainerID:   "",
				imageName:     "transeptorlabs/bundler:0.6.1-alpha.0",
				IsRunning:     false,
				Cmd:           []string{},
				ExposedPorts:  nil,
			},
			"geth": {
				containerName: "betsy-geth",
				ContainerID:   "",
				imageName:     "ethereum/client-go:latest",
				IsRunning:     false,
				Cmd: []string{
					"--dev",
					"--nodiscover",
					"--http",
					"--dev.gaslimit", "12000000",
					"--http.api", "eth,net,web3,debug",
					"--http.corsdomain", "*://localhost:*",
					"--http.vhosts", "*,localhost,host.docker.internal",
					"--http.addr", "0.0.0.0",
					"--networkid", "1337",
					"--verbosity", "1",
					"--maxpeers", "0",
					"--allow-insecure-unlock",
					"--rpc.allow-unprotected-txs",
				},
				ExposedPorts: nil,
			},
		},
		client: cli,
		ctx:    context.Background(),
	}, nil
}

// Close closes the Docker client
func (cm *ContainerManager) Close() error {
	return cm.client.Close()
}

func SmokeTestDockerAPI() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", image.PullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

// ListAllImages lists all images available in the Docker environment
func (cm *ContainerManager) ListRunningContainer() error {
	containers, err := cm.client.ContainerList(cm.ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		log.Info().Msgf("Container ID: %s\n", container.ID)
	}

	return nil
}

// PullRequiredImages checks if required images are available and pulls them if not
func (cm *ContainerManager) PullRequiredImages(requiredImages []string) (bool, error) {
	for _, requiredImage := range requiredImages {
		if _, ok := cm.supportedImages[requiredImage]; !ok {
			return false, fmt.Errorf("Image %s is not supported", requiredImage)
		}
	}

	requiredImageFoundCheck := make(map[string]bool)
	requiredImageNames := make([]string, 0)
	for _, requiredImage := range requiredImages {
		requiredImageFoundCheck[cm.supportedImages[requiredImage].imageName] = false
		requiredImageNames = append(requiredImageNames, cm.supportedImages[requiredImage].imageName)
	}
	log.Info().Msgf("Required images(pre-check): %v", requiredImageFoundCheck)

	localImages, err := cm.client.ImageList(cm.ctx, image.ListOptions{})
	if err != nil {
		return false, err
	}

	// Check if required images are found locally and update the check map with a caveat for the latest tag to ensure it is always pulled
	for _, image := range localImages {
		for _, requiredImageName := range requiredImageNames {
			if requiredImageName == image.RepoTags[0] && !strings.HasSuffix(requiredImageName, "latest") {
				requiredImageFoundCheck[requiredImageName] = true
			}
		}
	}
	log.Info().Msgf("Required images(post-check): %v", requiredImageFoundCheck)

	// Pull required images that are not found
	for imageName, found := range requiredImageFoundCheck {
		if !found {
			_, err := cm.doPullImage(imageName)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// doPullImage pulls a Docker image given its name
func (cm *ContainerManager) doPullImage(imageName string) (bool, error) {
	log.Info().Msgf("Attempting to pull image: %s", imageName)
	reader, err := cm.client.ImagePull(cm.ctx, imageName, image.PullOptions{})
	if err != nil {
		return false, err
	}
	io.Copy(os.Stdout, reader)

	log.Info().Msgf("Successfully pulled image: %s", imageName)
	return true, nil
}

// RunContainerInTheBackground runs a Docker container in the background given its image and host port to bind
func (cm *ContainerManager) RunContainerInTheBackground(image string, hostPort string) (bool, error) {
	imageFound, ok := cm.supportedImages[image]
	if !ok {
		return false, fmt.Errorf("Image %s is not supported", image)
	}

	constinerPort := hostPort + "/tcp"
	config := &container.Config{
		Image: imageFound.imageName,
		Cmd:   imageFound.Cmd,
		ExposedPorts: nat.PortSet{
			nat.Port(constinerPort): struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(constinerPort): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0", // setting to 0.0.0.0 means that the port is exposed on all network interfaces on host machine
					HostPort: hostPort,
				},
			},
		},
	}

	resp, err := cm.client.ContainerCreate(cm.ctx, config, hostConfig, nil, nil, imageFound.containerName)
	if err != nil {
		return false, err
	}

	if err := cm.client.ContainerStart(cm.ctx, resp.ID, container.StartOptions{}); err != nil {
		return false, err
	}

	// Update the container details
	log.Info().Msgf("Container ID successfully started: %s\n", resp.ID)
	cm.supportedImages[image] = ContainerDetails{
		imageName: imageFound.imageName,
		Cmd:       imageFound.Cmd,
		ExposedPorts: nat.PortSet{
			nat.Port(constinerPort): struct{}{},
		},
		ContainerID: resp.ID,
		IsRunning:   true,
	}

	return true, nil
}

// StopAndRemoveRunningContainers stops all running containers that are supported
func (cm *ContainerManager) StopAndRemoveRunningContainers() (bool, error) {
	for _, containerDetails := range cm.supportedImages {
		if containerDetails.IsRunning {
			log.Info().Msgf("Attempting to stop container %s", containerDetails.ContainerID)
			noWaitTimeout := 0 // to not wait for the container to exit gracefully

			if err := cm.client.ContainerStop(cm.ctx, containerDetails.ContainerID, container.StopOptions{Timeout: &noWaitTimeout}); err != nil {
				return false, err
			}
			log.Info().Msgf("Successfully stopped container %s", containerDetails.ContainerID)

			if err := cm.client.ContainerRemove(cm.ctx, containerDetails.ContainerID, container.RemoveOptions{}); err != nil {
				return false, err
			}
			log.Info().Msgf("Successfully removed container %s", containerDetails.ContainerID)
		}
	}

	return true, nil
}
