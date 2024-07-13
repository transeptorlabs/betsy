#!/bin/bash
EP_COMPILED_CONTRACT_FILE="$PWD/lib/account-abstraction-releases-v0.7/artifacts/contracts/core/EntryPoint.sol/EntryPoint.json"
EP_ABI_FILE="$PWD/precompiled-contracts/EntryPointV7.abi"
EP_BIN_FILE="$PWD/precompiled-contracts/EntryPointV7.bin"
EP_GO_BINDING_FILE="$PWD/contracts/entrypoint/entrypoint-v7.go"

SIMPLE_AF_COMPILED_CONTRACT_FILE="$PWD/lib/account-abstraction-releases-v0.7/artifacts/contracts/samples/SimpleAccountFactory.sol/SimpleAccountFactory.json"
SIMPLE_AF_ABI_FILE="$PWD/precompiled-contracts/SimpleAccountFactoryV7.abi"
SIMPLE_AF_BIN_FILE="$PWD/precompiled-contracts/SimpleAccountFactoryV7.bin"
SIMPLE_AF_GO_BINDING_FILE="$PWD/contracts/factory/simple-account-factory-v7.go"

GLOBAL_COUNTER_ABI_FILE="$PWD/precompiled-contracts/GlobalCounter.abi"
GLOBAL_COUNTER_BIN_FILE="$PWD/precompiled-contracts/GlobalCounter.bin"
GLOBAL_COUNTER_GO_BINDING_FILE="$PWD/contracts/examples/global-counter.go"

# check that jq is installed and exit if not
if ! [ -x "$(command -v jq)" ]; then
  echo "Error: jq is not installed. Please install to run script." >&2
  exit 1
fi

# Check that go is installed and exit if not
if ! [ -x "$(command -v go)" ]; then
  echo "Error: go is not installed. Please install to run script." >&2
  exit 1
fi

# Extract ABI and save to .abi file from the build artifacts and move to the precompiled-contracts shared location
# Set GOPATH for the script
export GOPATH=$HOME/go

# Ensure GOPATH/bin is in PATH
export PATH=$PATH:$GOPATH/bin

# Installing geth abigen tool
echo "Installing latest abigen..."
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# change to the contract directory
cd $PWD/lib/account-abstraction-releases-v0.7
if [ ! -d "node_modules" ]; then
  echo "Installing contract dependencies..."
  yarn
fi

# Compile contracts and check if the artifacts are generated
echo "Compiling account-abstraction-releases-v0.7 contracts..."
yarn compile

if [ ! -f $EP_COMPILED_CONTRACT_FILE ]; then
  echo "Failed to compile EntryPoint contract."
fi

if [ ! -f $SIMPLE_AF_COMPILED_CONTRACT_FILE ]; then
  echo "Failed to compile SimpleAccountFactory contract."
fi

echo "Extracting EntrypointV7 ABI and Bytecode..."
jq -r '.abi' $EP_COMPILED_CONTRACT_FILE > $EP_ABI_FILE
jq -r '.bytecode' $EP_COMPILED_CONTRACT_FILE > $EP_BIN_FILE

echo "EntrypointV7 ABI extracted to $EP_ABI_FILE"
echo "EntrypointV7 Bytecode extracted to $EP_BIN_FILE"

echo "Extracting SimpleAccountFactoryV7 ABI and Bytecode..."
jq -r '.abi' $SIMPLE_AF_COMPILED_CONTRACT_FILE > $SIMPLE_AF_ABI_FILE
jq -r '.bytecode' $SIMPLE_AF_COMPILED_CONTRACT_FILE > $SIMPLE_AF_BIN_FILE

echo "SimpleAccountFactoryV7 ABI extracted to $SIMPLE_AF_ABI_FILE"
echo "SimpleAccountFactoryV7 Bytecode extracted to $SIMPLE_AF_BIN_FILE"

# Generate the go bindings for the contracts
echo "Generating go bindings for EntrypointV7 contract..."
abigen --abi $EP_ABI_FILE --pkg entrypoint --type EntryPointV7 --bin $EP_BIN_FILE --out $EP_GO_BINDING_FILE

echo "Generating go bindings for SimpleAccountFactoryV7 contract..."
abigen --abi $SIMPLE_AF_ABI_FILE --pkg factory --type SimpleAccountFactoryV7 --bin $SIMPLE_AF_BIN_FILE --out $SIMPLE_AF_GO_BINDING_FILE

echo "Generating go bindings for GlobalCounter contract..."
abigen --abi $GLOBAL_COUNTER_ABI_FILE --pkg examples --type GlobalCounter --bin $GLOBAL_COUNTER_BIN_FILE --out $GLOBAL_COUNTER_GO_BINDING_FILE
