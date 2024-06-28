#!/bin/bash
# TODO: Move this script to golang to run as a sub command on besty

# Get the list of accounts
accounts_response=$(curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545)

# Extract the first account from the response
coinbase=$(echo $accounts_response | grep -o '"0x[^"]*"' | head -1)
coinbase=${coinbase//\"/}

echo "Coinbase account: $coinbase"

# Get the balance of the first account
balance_response=$(curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["'$coinbase'", "latest"],"id":1}' -H "Content-Type: application/json" http://localhost:8545)

# Extract the balance from the response
balance_hex=$(echo $balance_response | grep -o '"result":"0x[^"]*"' | awk -F'"' '{print $4}')

# Convert balance from hexadecimal to decimal
balance=$(printf "%d\n" $balance_hex)

echo "Balance: $balance wei"
