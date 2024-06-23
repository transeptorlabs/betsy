#  CLI Tool

This CLI tool allow you to manage ERC 4337 infrastructure using Docker containers.  An all in one tool to manage ERC 4337 infrastructure for local development and testing. The tool provides:
1. ETH client(i.e execution client)
   1. Forking mainnets and testnets
2. Default accounts with private keys
3. Predeployed entrypoint contract(V7)
4. ERC 4337 bundler client with 
5. ERC 4337 memepool/bundle explorer UI
6. ERC 4337 entrypoint contract UI(stake, unstake, deposit, withdraw)

## Command-Line Flags/Arguments
Use command-line flags to pass options like ports, network settings, etc., when starting containers.

## Error Handling
Implement robust error handling to manage potential issues such as Docker not being installed, Docker daemon not running, or containers failing to start.