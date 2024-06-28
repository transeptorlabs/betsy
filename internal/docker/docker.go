package docker

import (
	"fmt"
	"os"
	"strings"
	"time"

	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

const EthNodeReady = "ethNodeReady"

// ContainerManager manages containers
type ContainerManager struct {
	supportedImages map[string]ContainerDetails
	client          *client.Client
}

// ContainerDetails contains details of a container
type ContainerDetails struct {
	imageName     string
	containerName string
	ContainerID   string
	IsRunning     bool
	Cmd           []string
	ExposedPorts  nat.PortSet
	NodeType      string
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
				imageName:     "transeptorlabs/bundler:0.6.2-alpha.0", // Betsy Ross - https://github.com/transeptorlabs/transeptor-bundler/releases/tag/v0.6.2-alpha.0
				IsRunning:     false,
				Cmd:           []string{},
				ExposedPorts:  nil,
				NodeType:      "bundler",
			},
			"geth": {
				containerName: "betsy-geth",
				ContainerID:   "",
				imageName:     "ethereum/client-go:v1.14.5", // Bothros - https://github.com/ethereum/go-ethereum/releases/tag/v1.14.5
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
					"--verbosity", "2",
					"--maxpeers", "0",
					"--allow-insecure-unlock",
					"--rpc.allow-unprotected-txs",
				},
				ExposedPorts: nil,
				NodeType:     "eth",
			},
		},
		client: cli,
	}, nil
}

// Close closes the Docker client
func (cm *ContainerManager) Close() error {
	return cm.client.Close()
}

// ListAllImages lists all images available in the Docker environment
func (cm *ContainerManager) ListRunningContainer(ctx context.Context) error {
	containers, err := cm.client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		log.Info().Msgf("Container ID: %s\n", container.ID)
	}

	return nil
}

// PullRequiredImages checks if required images are available and pulls them if not
func (cm *ContainerManager) PullRequiredImages(ctx context.Context, requiredImages []string) (bool, error) {
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

	localImages, err := cm.client.ImageList(ctx, image.ListOptions{})
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
			_, err := cm.doPullImage(ctx, imageName)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// doPullImage pulls a Docker image given its name
func (cm *ContainerManager) doPullImage(ctx context.Context, imageName string) (bool, error) {
	log.Info().Msgf("Attempting to pull image: %s", imageName)
	reader, err := cm.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return false, err
	}
	io.Copy(os.Stdout, reader)

	log.Info().Msgf("Successfully pulled image: %s", imageName)
	return true, nil
}

// RunContainerInTheBackground runs a Docker container in the background given its image and host port to bind
func (cm *ContainerManager) RunContainerInTheBackground(ctx context.Context, image string, hostPort string) (bool, error) {
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
		// TODO: Use health check for bundlers and eth node
		// Healthcheck: &container.HealthConfig{
		// 	Test:     []string{"CMD", "curl", "-f", "http://localhost:" + hostPort},
		// 	Interval: 10 * time.Second,
		// 	Timeout:  5 * time.Second,
		// 	Retries:  5,
		// },
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

	resp, err := cm.client.ContainerCreate(ctx, config, hostConfig, nil, nil, imageFound.containerName)
	if err != nil {
		return false, err
	}

	if err := cm.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
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

	// Update EthNodeReady channel and signal that eth is ready by closing the channel
	if imageFound.NodeType == "eth" {
		log.Info().Msg("Waiting for Eth node container to become ready...")

		for {
			containerJSON, err := cm.client.ContainerInspect(ctx, resp.ID)
			if err != nil {
				return false, err
			}

			log.Info().Msgf("Eth node container is ready status: %+v", containerJSON.State.Status)
			if containerJSON.State.Status == "running" {
				break
			}
			time.Sleep(3 * time.Second)
		}

		if readyChan, ok := ctx.Value(EthNodeReady).(chan struct{}); ok {
			close(readyChan)
		}
	}

	return true, nil
}

// StopAndRemoveRunningContainers stops all running containers that are supported
func (cm *ContainerManager) StopAndRemoveRunningContainers(ctx context.Context) (bool, error) {
	for _, containerDetails := range cm.supportedImages {
		if containerDetails.IsRunning {
			log.Info().Msgf("Attempting to stop container %s", containerDetails.ContainerID)
			noWaitTimeout := 0

			if err := cm.client.ContainerStop(ctx, containerDetails.ContainerID, container.StopOptions{Timeout: &noWaitTimeout}); err != nil {
				return false, err
			}
			log.Info().Msgf("Successfully stopped container %s", containerDetails.ContainerID)

			if err := cm.client.ContainerRemove(ctx, containerDetails.ContainerID, container.RemoveOptions{}); err != nil {
				return false, err
			}
			log.Info().Msgf("Successfully removed container %s", containerDetails.ContainerID)
		}
	}

	return true, nil
}
