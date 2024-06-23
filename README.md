# 4337 In a Box

[![API Reference](
https://pkg.go.dev/badge/github.com/transeptorlabs/betsy
)](https://pkg.go.dev/github.com/transeptorlabs/betsy)
[![Go Report Card](https://goreportcard.com/badge/github.com/transeptorlabs/betsy)](https://goreportcard.com/report/github.com/transeptorlabs/betsy)

An all in one tool to manage ERC 4337 infrastructure for local development and testing. The tool provides:
1. ETH client(i.e execution client)
   - Forking evm `mainnets` and `testnets`
2. Default accounts with private keys
3. Predeployed entrypoint contract(V7)
4. ERC 4337 bundler client with 
5. ERC 4337 memepool/bundle explorer UI
6. ERC 4337 entrypoint contract UI(stake, unstake, deposit, withdraw)

**Supported ERC 4337 bundlers**
- [x] [Transeptor](https://github.com/transeptorlabs/transeptor-bundler)
- **Other bundlers coming soon**

## Installation

**Requirements**:
1. [Go - >= v1.22.4](https://go.dev/doc/install)
2. [Docker](https://docs.docker.com/engine/install)

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
docker run -it --rm 4337-in-a-box:v-local --help
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
1. Fork the repository
2. Create your feature branch (git checkout -b feature/fooBar)
3. Commit your changes (git commit -m 'Add some fooBar')
4. Push to the branch (git push origin feature/fooBar)
5. Open a Pull Request
