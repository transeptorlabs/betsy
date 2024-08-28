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
	"github.com/transeptorlabs/betsy/wallet"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

const EthNodeReady = "ethNodeReady"
const EthNodePortPlaceHolder = "$ETH_PORT"

const BundlerNodeWalletDetails = "bundlerNodeWalletDetails"
const BundlerNodeEPAddressPlaceHolder = "$ENTRYPOINT_ADDRESS"
const BundlerNodeBeneficiaryAddressPlaceHolder = "$BENEFICIARY"
const BundlerNodeMnemonicPlaceHolder = "$MNEMONIC"

// ContainerManager manages containers
type ContainerManager struct {
	supportedImages      map[string]*ContainerDetails
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

// NewContainerManager creates a new container manager
func NewContainerManager() (*ContainerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &ContainerManager{
		supportedImages: map[string]*ContainerDetails{
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
					"--autoBundleInterval", "10000", // 10 secs
					"--network", "http://host.docker.internal:" + EthNodePortPlaceHolder,
				},
				Env: []string{
					"TRANSEPTOR_MNEMONIC=" + BundlerNodeMnemonicPlaceHolder,
					"TRANSEPTOR_BENEFICIARY=" + BundlerNodeBeneficiaryAddressPlaceHolder,
					"TRANSEPTOR_ENTRYPOINT_ADDRESS=" + BundlerNodeEPAddressPlaceHolder,
				},
				ExposedPorts: nil,
				NodeType:     "bundler",
			},
			"aabundler": {
				containerName: "betsy-aabundler",
				ContainerID:   "",
				imageName:     "accountabstraction/bundler:0.7.0",
				IsRunning:     false,
				Cmd: []string{
					"--network", "http://host.docker.internal:" + EthNodePortPlaceHolder,
					"--entryPoint", BundlerNodeEPAddressPlaceHolder,
					"--beneficiary", BundlerNodeBeneficiaryAddressPlaceHolder,
				},
				Env:          []string{},
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
					"--dev.gaslimit", "30000000",
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
		log.Debug().Msgf("Container ID: %s\n", container.ID)
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

func (cm *ContainerManager) doReplaceBundlerPlaceHolders(containerDetails *ContainerDetails, bwd wallet.BundlerWalletDetails) {
	placeholders := map[string]string{
		BundlerNodeEPAddressPlaceHolder:          bwd.EntryPointAddress.Hex(),
		BundlerNodeBeneficiaryAddressPlaceHolder: bwd.Beneficiary.Hex(),
		BundlerNodeMnemonicPlaceHolder:           bwd.Mnemonic,
		EthNodePortPlaceHolder:                   cm.EthNodePort,
	}

	// Replace placeholders in Cmd slice
	for placeholder, value := range placeholders {
		replacePlaceHolderInSlice(&containerDetails.Cmd, placeholder, value)
	}

	// Replace placeholders in Env slice

	for placeholder, value := range placeholders {
		replacePlaceHolderInSlice(&containerDetails.Env, placeholder, value)
	}
}

// replacePlaceHolderInSlice finds and replaces a placeholder in a slice of strings.
func replacePlaceHolderInSlice(slice *[]string, placeholder, replacement string) {
	for i, item := range *slice {
		if strings.Contains(item, placeholder) {
			(*slice)[i] = strings.Replace(item, placeholder, replacement, 1)
		}
	}
}

// RunContainerInTheBackground runs a Docker container in the background given its image and host port to bind
func (cm *ContainerManager) RunContainerInTheBackground(ctx context.Context, image string, hostPort string) (bool, error) {
	imageFound, ok := cm.supportedImages[image]
	if !ok {
		return false, fmt.Errorf("Image %s is not supported", image)
	}

	// Update bundler node cmd with ethNode port
	if imageFound.NodeType == "bundler" {
		bwd := ctx.Value(BundlerNodeWalletDetails).(wallet.BundlerWalletDetails)
		cm.doReplaceBundlerPlaceHolders(imageFound, bwd)
	}

	// Create and start the container
	containerPort := hostPort + "/tcp"
	config := &container.Config{
		Image: imageFound.imageName,
		Cmd:   imageFound.Cmd,
		Env:   imageFound.Env,
		ExposedPorts: nat.PortSet{
			nat.Port(containerPort): struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(containerPort): []nat.PortBinding{
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
	log.Debug().Msgf("%s Container ID successfully started: %s\n", image, resp.ID)
	imageFound.ExposedPorts = nat.PortSet{
		nat.Port(containerPort): struct{}{},
	}
	imageFound.ContainerID = resp.ID
	imageFound.IsRunning = true

	// Update EthNodeReady channel and signal that eth is ready by closing the channel
	if imageFound.NodeType == "eth" {
		log.Info().Msg("Waiting for Eth node container to become ready...")
		for {
			containerJSON, err := cm.client.ContainerInspect(ctx, resp.ID)
			if err != nil {
				return false, err
			}

			log.Debug().Msgf("Checking Eth node container ready status: %+v", containerJSON.State.Status)
			if containerJSON.State.Status == "running" {
				break
			}
			time.Sleep(3 * time.Second)
		}

		time.Sleep(3 * time.Second)
		log.Debug().Msgf("Attempting to find eth.coinbase keystore file at /tmp on container: %s", resp.ID)
		coinbaseKeystoreFile, err := findCoinbaseKeystoreFile(resp.ID, "tmp")
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

// StopAndRemoveRunningContainers stops all running containers that are supported
func (cm *ContainerManager) StopAndRemoveRunningContainers(ctx context.Context) (bool, error) {
	for _, containerDetails := range cm.supportedImages {
		if containerDetails.IsRunning {
			log.Debug().Msgf("Attempting to stop container %s", containerDetails.ContainerID)
			noWaitTimeout := 0

			if err := cm.client.ContainerStop(ctx, containerDetails.ContainerID, container.StopOptions{Timeout: &noWaitTimeout}); err != nil {
				return false, err
			}
			log.Debug().Msgf("Successfully stopped container %s", containerDetails.ContainerID)

			if err := cm.client.ContainerRemove(ctx, containerDetails.ContainerID, container.RemoveOptions{}); err != nil {
				return false, err
			}
			log.Debug().Msgf("Successfully removed container %s", containerDetails.ContainerID)
		}
	}

	return true, nil
}

// findCoinbaseKeystoreFile executes the ls command to find keystore file for the temporary pre-allocated developer account available and unlocked as eth.coinbase(using docker exec)
func findCoinbaseKeystoreFile(containerID string, dir string) (string, error) {
	cmd := exec.Command("docker", "exec", containerID, "ls", dir)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Recursive find
	log.Debug().Msgf("Files in directory: %s\n%s", dir, string(output))
	fileList := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range fileList {
		if strings.Contains(file, "UTC") {
			foundPath := strings.TrimSuffix(dir, "/") + "/" + file
			log.Debug().Msgf("Found keystore file path: %s", foundPath)

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
			log.Debug().Msgf("Not found, searching deeper in directory: %s", newDir)
			foundPath, err := findCoinbaseKeystoreFile(containerID, newDir)
			if err == nil {
				return foundPath, nil
			}
		}
	}

	// Keystore file not found
	log.Warn().Msgf("Keystore file not found in directory: %s", dir)
	return "", fmt.Errorf("keystore file not found in directory: %s", dir)
}

// copyFileFromContainer copies a file from a Docker container to the local filesystem
func copyFileFromContainer(containerID, filePath string, destDir string) error {
	destLocalFilePath := filepath.Join(destDir, filepath.Base(filePath))

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
	destFile, err := os.Create(destLocalFilePath)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating file %s", destLocalFilePath)
		return err
	}
	defer destFile.Close()

	// Write the file contents to the local file
	_, err = destFile.Write(output)
	if err != nil {
		log.Error().Err(err).Msgf("Error writing to file %s", destLocalFilePath)
		return err
	}

	log.Debug().Msgf("Copied file %s from container %s to %s", filePath, containerID, destLocalFilePath)
	return nil
}
