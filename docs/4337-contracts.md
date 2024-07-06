# ERC 4337 Contracts

We use Geth [abigen](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings), a tool that converts Ethereum ABI definitions into easy-to-use, type-safe Go packages to interact programmatically with ERC 4337 Contracts. The contract binding is used to deploy the contracts on Betsy's start-up.

**Supported contracts:**
- [Entytpoint release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/core/EntryPoint.sol)
- [SimpleAccountFactory release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/samples/SimpleAccountFactory.sol)
- [SimpleAccount release v7](https://github.com/eth-infinitism/account-abstraction/blob/releases/v0.7/contracts/samples/SimpleAccount.sol)

To generate the contract bindings, run the following command:
```bash
make gen-contract-binding-aa
```