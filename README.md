# Betsy

[![API Reference](
https://pkg.go.dev/badge/github.com/transeptorlabs/betsy
)](https://pkg.go.dev/github.com/transeptorlabs/betsy)
[![Go Report Card](https://goreportcard.com/badge/github.com/transeptorlabs/betsy)](https://goreportcard.com/report/github.com/transeptorlabs/betsy)
![build status](https://github.com/transeptorlabs/betsy/actions/workflows/build.yml/badge.svg?branch=main)
[![Docker Pulls](https://img.shields.io/docker/pulls/transeptorlabs/betsy)](https://img.shields.io/docker/pulls/transeptorlabs/betsy)


An all in one cli tool to manage ERC 4337 infrastructure for local development and testing. The tool provides:
1. Ethereum execution client with support to fork EVM mainnet and testnet
2. Uses an ephemeral in-memory execution client that is completely destroyed and starts a fresh instance during each Betsy run
3. Prefunded default accounts with private keys
4. Predeployed entrypoint point [releases/v0.7](https://github.com/eth-infinitism/account-abstraction/tree/releases/v0.7) contract - 
5. Pick from multiple ERC 4337 Bundler clients
6. Realtime ERC 4337 userOp mempool explorer UI
7. Realtime ERC 4337 Bundler bundle explorer UI

**Supported ERC 4337 bundlers**
- [x] [Transeptor](https://github.com/transeptorlabs/transeptor-bundler)
- **Other bundlers coming soon**

## Installation

**Requirements**:
1. [Go - >= v1.22.4](https://go.dev/doc/install)
2. [Docker - >= 20.10.17](https://docs.docker.com/engine/install)

### Build from the source

To build the project from the source code, run:
```shell
make build-source
```

### Build Docker image

To build a Docker image for the project, run:

```shell
make build-docker
```

Start the container with the following command:  
```shell
docker run -it --rm besty:v-local --help
```

##  Development

### Running the application

Run the following command to start the application:
```shell
make run-app
```

### Running tests

To run the tests, execute the following command:
```shell
make test
```

or to run the tests with coverage:
```shell
make test-coverage
```

##  Contributing

If you would like to contribute, please follow these guidelines [here](https://github.com/transeptorlabs/betsy/blob/main/CONTRIBUTING.md).

## License

Licensed under the [MIT](https://github.com/transeptorlabs/betsy/blob/main/LICENSE).
