#  CLI Tool Details

## Docker Engine API

The [Docker Engine API](https://docs.docker.com/engine/api/) manages all erc4337 environment containers. Under the hood, Betsy plugs into the API to create, start, stop, and remove containers, abstracting the complexity of managing Docker containers.

## Command-Line Flags/Arguments

When starting containers, use command-line flags to pass options like ports, network settings, etc.

## Error Handling

Implement robust error handling to manage potential issues such as Docker not being installed, Docker daemon not running, or containers failing to start.

## Project structure

This project structure adheres to standard Go project layouts:
```txt
4337-in-a-box/
├── cmd/            # Main applications of the project
│   └── 4337-in-a-box/    # Application entry point
│       └── main.go
├── internal/       # Private application and library code
├── precompiled-contracts/      # Local precompiled contracts
├── version/        # Version information
├── go.mod          # Module definition
├── go.sum          # Dependencies checksum
└── README.md       # Project documentation
```

## References
- Docker Engine API: https://docs.docker.com/engine/api/
- The Moby Project: https://pkg.go.dev/github.com/docker/docker#section-readme
  - Go client for the Docker Engine API: https://pkg.go.dev/github.com/docker/docker/client#section-readme