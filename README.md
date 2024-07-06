# Betsy

[![API Reference](
https://pkg.go.dev/badge/github.com/transeptorlabs/betsy
)](https://pkg.go.dev/github.com/transeptorlabs/betsy)
[![Go Report Card](https://goreportcard.com/badge/github.com/transeptorlabs/betsy)](https://goreportcard.com/report/github.com/transeptorlabs/betsy)
![build status](https://github.com/transeptorlabs/betsy/actions/workflows/build.yml/badge.svg?branch=main)

An all-in-one CLI tool to manage ERC 4337 infrastructure for local development and testing. 

✨ **Features include:**
1. Uses an ephemeral in-memory Ethereum execution client that is completely destroyed and starts a fresh instance during each Betsy run
2. Pre-funded accounts: Default pre-funded accounts with private keys
3. Pre-deployed contract: Type-safe Go binding for Account abstraction contracts
   - [EntryPoint release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/core/EntryPoint.sol)
   - [SimpleAccountFactory release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/samples/SimpleAccountFactory.sol)
   - [SimpleAccount release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/samples/SimpleAccount.sol)
4. ERC 4337 Bundler clients
    - [x] [Transeptor](https://github.com/transeptorlabs/transeptor-bundler)
5. Realtime ERC 4337 userOp mempool explorer UI. Visualize the userOp mempool in real time.
6. Realtime ERC 4337 Bundler bundle explorer UI. Visualize the bundler bundle production in real time.

🚧 **Coming soon:**
1. Supported ERC 4337 bundlers**
   - [x] [Transeptor](https://github.com/transeptorlabs/transeptor-bundler)
   - [ ] Other bundlers (e.g. [Aabundler](https://github.com/eth-infinitism/bundler), [Okbund](https://github.com/okx/okbund) etc.)
2. Ethereum execution client forks for EVM mainnet and testnet. Run a fork of the Ethereum mainnet or testnet to test your AA smart contracts in a real-world environment.
3. Manage Entrypoint deposits, withdrawals, and stakes on local Entrypoint contract.

## Installation

For more information on installing Betsy, see the [Installation](./docs/installation.md) guide.

##  Development

Information on how to set a development environment for Betsy.

### Branches

The `main` branch acts as the development branch and is the repository's default branch. The main branch build will be marked as `unstable` in the version.
- Betsy's latest `stable` version can be found on branch `release/x.y.z`.

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

Or to run the tests with coverage:
```shell
make test-coverage
```

##  Contributing

If you would like to contribute, please follow these guidelines [here](https://github.com/transeptorlabs/betsy/blob/main/CONTRIBUTING.md).

## License

Licensed under the [MIT](https://github.com/transeptorlabs/betsy/blob/main/LICENSE).
