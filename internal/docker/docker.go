package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
const EthNodePortPlaceHolder = "$ETH_PORT"

// ContainerManager manages containers
type ContainerManager struct {
	supportedImages      map[string]ContainerDetails
	client               *client.Client
	EthNodePort          string
	CoinbaseKeystoreFile string
}

// ContainerDetails contains details of a container
type ContainerDetails struct {
	imageName     string
	containerName string
	ContainerID   string
	IsRunning     bool
	Cmd           []string
	Env           []string
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
				Cmd: []string{
					"--txMode", "base",
					"--unsafe",
					"--httpApi", "web3,eth,debug",
					"--auto",
					"--autoBundleInterval", "12000",
					"--network", "http://host.docker.internal:" + EthNodePortPlaceHolder,
				},
				Env: []string{
					"TRANSEPTOR_MNEMONIC=test test test test test test test test test test test junk",
					"TRANSEPTOR_BENEFICIARY=0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
				},
				ExposedPorts: nil,
				NodeType:     "bundler",
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
				Env:          []string{},
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

// IsDockerInstalled checks if docker is installed
func (cm *ContainerManager) IsDockerInstalled() bool {
	// Command to check if Docker is installed
	cmd := exec.Command("docker", "version")

	// Try running the command
	err := cmd.Run()

	// If there's an error, Docker is not installed
	if err != nil {
		return false
	}

	return true
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

	// Update bundler node cmd with ethnode port
	if imageFound.NodeType == "bundler" {
		foundIndex := 0
		for index, item := range imageFound.Cmd {
			if strings.HasSuffix(item, "ETH_PORT") {
				foundIndex = index
				break
			}
		}
		imageFound.Cmd[foundIndex] = strings.Replace(imageFound.Cmd[foundIndex], EthNodePortPlaceHolder, cm.EthNodePort, 1)
	}

	constinerPort := hostPort + "/tcp"
	config := &container.Config{
		Image: imageFound.imageName,
		Cmd:   imageFound.Cmd,
		Env:   imageFound.Env,
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

			log.Info().Msgf("Checking Eth node container ready status: %+v", containerJSON.State.Status)
			if containerJSON.State.Status == "running" {
				break
			}
			time.Sleep(3 * time.Second)
		}

		log.Info().Msgf("Attempting to find eth.coinbase keystore file at /tmp on container: %s", resp.ID)
		coinbaseKeystoreFile, err := findCoinbaseKeystoreFileNative(resp.ID, "tmp")
		if err != nil {
			return false, err
		}

		cm.EthNodePort = hostPort
		cm.CoinbaseKeystoreFile = coinbaseKeystoreFile

		if readyChan, ok := ctx.Value(EthNodeReady).(chan struct{}); ok {
			close(readyChan)
		}
	}

	return true, nil
}

// TODO: Fix the EOF error when runing the ls command on the container
// findCoinbaseKeystoreFile executes the ls command to find keystore file for the temporary pre-allocated developer account available and unlocked as eth.coinbase(using docker api client)
func findCoinbaseKeystoreFile(ctx context.Context, client *client.Client, containerID string, dir string) (string, error) {
	execConfig := container.ExecOptions{
		Cmd:          []string{"ls", dir},
		AttachStdout: true,
		AttachStderr: true,
	}
	execIDResp, err := client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", err
	}

	respAttach, err := client.ContainerExecAttach(ctx, execIDResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", err
	}
	defer respAttach.Close()

	output := make([]byte, 1024)
	n, err := respAttach.Reader.Read(output)
	if err != nil {
		fmt.Println("Got a error when listing:", err)
		return "", err
	}

	fmt.Printf("Files in /tmp:\n%s\n", string(output[:n]))
	return string(output[:n]), nil
}

// findCoinbaseKeystoreFileNative executes the ls command to find keystore file for the temporary pre-allocated developer account available and unlocked as eth.coinbase(using docker exec)
func findCoinbaseKeystoreFileNative(containerID string, dir string) (string, error) {
	cmd := exec.Command("docker", "exec", containerID, "ls", dir)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Print the output for debugging purposes
	fmt.Println("Files in directory:", dir)
	fmt.Println(string(output))

	// Recursive find
	fileList := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range fileList {
		if strings.Contains(file, "UTC") {
			foundPath := strings.TrimSuffix(dir, "/") + "/" + file
			log.Info().Msgf("Found keystore file path: %s", foundPath)

			// Copy the found file to local ./wallet/tmp/coinbase directory
			if err := copyFileFromContainer(containerID, foundPath, "./wallet/tmp/coinbase"); err != nil {
				log.Error().Err(err).Msg("Error copying file from container")
				return "", err
			}

			return file, nil
		}
	}

	// If not found, recursively search deeper
	for _, file := range fileList {
		if !strings.Contains(file, ".") { // Assuming files without dots are directories
			newDir := strings.TrimSuffix(dir, "/") + "/" + file
			log.Info().Msgf("Not found, searching deeper in directory: %s", newDir)
			foundPath, err := findCoinbaseKeystoreFileNative(containerID, newDir)
			if err == nil {
				return foundPath, nil
			}
		}
	}

	// Keystore file not found
	log.Warn().Msgf("Keystore file not found in directory: %s", dir)
	return "", fmt.Errorf("keystore file not found in directory: %s", dir)
}

func copyFileFromContainer(containerID, filePath string, destDir string) error {
	// Determine the destination path on the host machine
	destPath := filepath.Join(destDir, filepath.Base(filePath))

	// Ensure the destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Error().Err(err).Msgf("Error creating directory %s", destDir)
		return err
	}

	// Open a reader to read the file contents from the Docker container
	cmd := exec.Command("docker", "exec", containerID, "cat", filePath)
	output, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msgf("Error reading file %s from container %s", filePath, containerID)
		return err
	}

	// Create a new file in local directory
	destFile, err := os.Create(destPath)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating file %s", destPath)
		return err
	}
	defer destFile.Close()

	// Write the file contents to the local file
	_, err = destFile.Write(output)
	if err != nil {
		log.Error().Err(err).Msgf("Error writing to file %s", destPath)
		return err
	}

	log.Info().Msgf("Copied file %s from container %s to %s", filePath, containerID, destPath)
	return nil
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
