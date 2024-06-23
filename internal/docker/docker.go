package docker

import (
	"fmt"
	"os"

	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// ContainerManger manages containers
type ContainerManger struct {
	supportedImages map[string]ContainerDetails
}

// ContainerDetails contains details of a container
type ContainerDetails struct {
	imageName   string
	ContainerID string
	IsRunning   bool
	Cmd         []string
	Port        string
}

// NewContainerManager creates a new container manager
func NewContainerManager() *ContainerManger {
	return &ContainerManger{
		supportedImages: map[string]ContainerDetails{
			"transeptor": {
				imageName:   "transeptorlabs/bundler:latest",
				ContainerID: "",
				IsRunning:   false,
				Port:        "",
			},
			"geth": {
				imageName:   "ethereum/client-go:latest",
				ContainerID: "",
				IsRunning:   false,
				Port:        "",
				Cmd: []string{
					"--verbosity", "1",
					"--http.vhosts", "'*,localhost,host.docker.internal'",
					"--http",
					"--http.api", "eth,net,web3,debug",
					"--http.corsdomain", "'*'",
					"--http.addr", "0.0.0.0",
					"--nodiscover", "--maxpeers", "0", "--mine",
					"--networkid", "1337",
					"--dev",
					"--allow-insecure-unlock",
					"--rpc.allow-unprotected-txs",
					"--dev.gaslimit", "12000000",
				},
			},
		},
	}
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
func (cm *ContainerManger) ListRunningContainer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

// PullRequiredImages checks if required images are available and pulls them if not
func (cm *ContainerManger) PullRequiredImages(requiredImages []string) (bool, error) {
	for _, requiredImage := range requiredImages {
		if _, ok := cm.supportedImages[requiredImage]; !ok {
			return false, fmt.Errorf("Image %s is not supported", requiredImage)
		}
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return false, err
	}

	foundImageName := make(map[string]bool)
	for _, requiredImage := range requiredImages {
		foundImageName[requiredImage] = false
	}

	// Check if required images are found
	for _, image := range images {
		for _, requiredImage := range requiredImages {
			fmt.Println(image.RepoTags)
			if requiredImage == image.RepoTags[0] {
				fmt.Println("Image found:", image.RepoTags[0])
				foundImageName[requiredImage] = true
			}
		}
	}

	// Pull images that are not found
	for imageName, found := range foundImageName {
		if !found {
			fmt.Println("Image not found, pulling:", imageName)
			_, err := cm.doPullImage(imageName)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

// doPullImage pulls a Docker image given its name
func (cm *ContainerManger) doPullImage(imageName string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return false, err
	}
	io.Copy(os.Stdout, reader)

	return true, nil
}

// RunContainerInTheBackground runs a Docker container in the background
func (cm *ContainerManger) RunContainerInTheBackground(image string, port string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}
	defer cli.Close()

	imageFound, ok := cm.supportedImages[image]
	if !ok {
		return false, fmt.Errorf("Image %s is not supported", image)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:    imageFound.imageName,
		Hostname: "localhost:" + port,
	}, nil, nil, nil, "")
	if err != nil {
		return false, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return false, err
	}

	cm.supportedImages[image] = ContainerDetails{
		imageName:   imageFound.imageName,
		ContainerID: resp.ID,
		IsRunning:   true,
		Port:        port,
	}
	fmt.Println(resp.ID, cm.supportedImages[image], "started successfully.")

	return true, nil
}

// StopRunningContainers stops all running containers that are supported
func (cm *ContainerManger) StopRunningContainers() (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}
	defer cli.Close()

	for _, containerDetails := range cm.supportedImages {
		if containerDetails.IsRunning {
			fmt.Print("Stopping container ", containerDetails.ContainerID, "... ")
			noWaitTimeout := 0 // to not wait for the container to exit gracefully
			if err := cli.ContainerStop(ctx, containerDetails.ContainerID, container.StopOptions{Timeout: &noWaitTimeout}); err != nil {
				return false, err
			}
			fmt.Println("Successfully stopped.")
		}
	}

	return true, nil
}
