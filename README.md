# 4337 In a Box

## Project structure

This project structure adheres to standard Go project layouts:
```txt
4337-in-a-box/
├── cmd/            # Main applications of the project
│   └── 4337-in-a-box/    # Application entry point
│       └── main.go
├── internal/       # Private application and library code
├── contracts/      # Forge smart contracts
├── version/        # Version information
├── scripts/        # Local scripts
├── go.mod          # Module definition
├── go.sum          # Dependencies checksum
└── README.md       # Project documentation
```

## Requirements

1. [Go - >= v1.22.4](https://go.dev/doc/install)
2. [Docker](https://docs.docker.com/engine/install)

## Supported 4337 bundlers
 - [x] Transeptor
  
**Other bundlers coming soon**

## Installation

### Clone the Repository
```shell
git clone https://github.com/transeptorlabs/4337-in-a-box.git
cd 4337-in-a-box
```

### Initialize Go Module
```shell
go mod tidy
```

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

### Contracts
1. `git submodule update --init`
2. `cd contracts`
3. `forge compile`

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

### Local Scripts
Local scripts are stored in the scripts directory. You can run these scripts as needed.

Start local eth node:
```shell
make eth
```

##  Contributing
1. Fork the repository
2. Create your feature branch (git checkout -b feature/fooBar)
3. Commit your changes (git commit -m 'Add some fooBar')
4. Push to the branch (git push origin feature/fooBar)
5. Open a Pull Request
