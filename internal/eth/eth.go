package eth

import (
	"fmt"
	"os"
	"os/exec"
)

var gethContainerName = "geth-4337-in-a-box"

func CheckGethNodeRunning() bool {
	return false
}

func StartGethNode() (bool, error) {
	// Set default values for GETH_PORT and NETWORK_ID
	gethPort := "8545"
	networkID := "1337"

	// Check if environment variables are set and use them if available
	if port, ok := os.LookupEnv("GETH_PORT"); ok {
		gethPort = port
	}
	if id, ok := os.LookupEnv("NETWORK_ID"); ok {
		networkID = id
	}

	// Example: Start Geth node in Docker container with specified options
	cmd := exec.Command("docker", "run", "-d", "--name", gethContainerName,
		"-p", fmt.Sprintf("%s:%s", gethPort, gethPort),
		"ethereum/client-go:latest",
		"--verbosity", "1",
		"--http.vhosts", "'*,localhost,host.docker.internal'",
		"--http",
		"--http.api", "eth,net,web3,debug",
		"--http.corsdomain", "'*'",
		"--http.addr", "0.0.0.0",
		"--nodiscover", "--maxpeers", "0", "--mine",
		"--networkid", networkID,
		"--dev",
		"--allow-insecure-unlock",
		"--rpc.allow-unprotected-txs",
		"--dev.gaslimit", "12000000")

	err := cmd.Run()
	if err != nil {
		return false, err
	}
	fmt.Println("Geth node started successfully.")
	return true, nil
}

func StopGethNode() {
	// Example: Stop Geth node Docker container
	cmd := exec.Command("docker", "stop", gethContainerName)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error stopping Geth node:", err)
		os.Exit(1)
	}
	fmt.Println("Geth node stopped successfully.")
}
