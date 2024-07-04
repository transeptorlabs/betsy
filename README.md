# Betsy

[![API Reference](
https://pkg.go.dev/badge/github.com/transeptorlabs/betsy
)](https://pkg.go.dev/github.com/transeptorlabs/betsy)
[![Go Report Card](https://goreportcard.com/badge/github.com/transeptorlabs/betsy)](https://goreportcard.com/report/github.com/transeptorlabs/betsy)
![build status](https://github.com/transeptorlabs/betsy/actions/workflows/build.yml/badge.svg?branch=main)
[![Docker Pulls](https://img.shields.io/docker/pulls/transeptorlabs/betsy)](https://img.shields.io/docker/pulls/transeptorlabs/betsy)


An all in one cli tool to manage ERC 4337 infrastructure for local development and testing. 

**Features include:**
1. Uses an ephemeral in-memory Ethereum execution client that is completely destroyed and starts a fresh instance during each Betsy run
2. Pre-funded accounts: Default pre-funded accounts with private keys
3. Pre-deployed contract: Type-safe Go binding for Account abstraction contracts
   - [EntryPoint release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/core/EntryPoint.sol)
   - [SimpleAccountFactory release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/samples/SimpleAccountFactory.sol)
   - [SimpleAccount release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/samples/SimpleAccount.sol)
4. ERC 4337 Bundler clients
    -  [x] [Transeptor](https://github.com/transeptorlabs/transeptor-bundler)
5. Realtime ERC 4337 userOp mempool explorer UI. Visualize the userOp mempool in real-time.
6. Realtime ERC 4337 Bundler bundle explorer UI. Visualize the bundler bundle production in real-time.

**Coming soon:**
1. Supported ERC 4337 bundlers**
   - [x] [Transeptor](https://github.com/transeptorlabs/transeptor-bundler)
   - [ ] Other bundlers (e.g. Flashbots, ArcherDAO, etc.)
2. Ethereum execution client forks for EVM mainnet and testnet. Run your own fork of the Ethereum mainnet or testnet to test your AA smart contracts in a real-world environment.
3. Manage Entrypoint deposits, withdrawals, and stakes on local Entrypoint contract.

## Versioning

Betsy follows [Semantic Versioning](https://semver.org/) for versioning releases. Each release can be found on the repository as a branch with the version number `/release/x.y.z.` along with a release tag with the version number `vx.y.z`.

### Branches

The `main` branch is the default branch for the repository and acts as the development branch. The main brach is `unstable` and should be for those who want to run the latest version of Besty to test new features and bug fixes.

The latest `stable` version of Besty can be found on branch `release/x.y.z`. The stable branch is for **those who want to run the latest stable version of Besty.

If you are unsure which version of Besty you are currently running, you can check the version by running the command `betsy --version`. You should see the version number, commit hash, and commit date for the latest stable/unstable version, in the following format.

For unstable version(development):
```shell    
betsy version x.y.z-unstable (abcabcabcabc yyyy-mm-dd)
```

For stable version:
```shell
betsy version x.y.z-stable (abcabcabcabc yyyy-mm-dd)
```

## Installation

**Requirements**:
1. [Go - >= v1.22.4](https://go.dev/doc/install)
2. [Docker](https://docs.docker.com/engine/install)

### Build from the source

#### Linux and Mac

:::important Betsy with default to unstable builds(development) If you want to use a stable build; please checkout a `release/x.y.z` branch before running `make besty`. :::

For UNIX-like operating systems you can clone the [Besty](https://github.com/transeptorlabs/betsy) repository and create a **temporary** build using the command `make besty`. This method of building requires Go(>= 1.22.4) and Docker to be installed on your system.

```shell
git clone https://github.com/transeptorlabs/betsy.git
cd besty
make besty
```

Running the command above results in the creation of a standalone executable file in the `betsy/bin` directory and does not require any dependencies to run. You can run the executable file using the command `./bin/betsy --help`. Or you can move the executable file and run from another directory.

To update the the latest version of Besty, you can:
1. Stop the cli(If it is running)
2. Navigate to the Besty directory and 
3. Pull the latest version of the source code from Besty Github repository 
4. Build and restart the cli

```shell
cd betsy
git pull
make betsy
```

##  Development

### Running the CLI

Run the following command to start the CLI:
```shell
make run-cli
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
