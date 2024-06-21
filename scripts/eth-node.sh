#!/bin/bash

GETH_PORT=8545
TRANSEPTOR_PORT=4337
NETWORK_ID=1337
TRANSEPTOR_ENTRYPOINT_ADDRESS_V7=""

start_geth() {
  # first 3 accounts default hardhat accounts
  DEFAULT_ADDRESS_1="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266" # bundler signer
  DEFAULT_ADDRESS_2="0x70997970C51812dc3A010C7d01b50e0d17dc79C8" # local e2e runner
  DEFAULT_ADDRESS_3="0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC" # bundler beneficiary

  if docker ps -a | grep -q geth-transeptor; then
    echo -e "Removing existing geth container\n"
    docker rm -f geth-transeptor
  fi

  echo -e "Starting local eth node at http://localhost:$GETH_PORT on network $NETWORK_ID\n"
  geth_container_id=$(docker run -d --name geth-transeptor -p $GETH_PORT:$GETH_PORT ethereum/client-go:latest \
    --verbosity 1 \
    --http.vhosts '*,localhost,host.docker.internal' \
    --http \
    --http.api eth,net,web3,debug \
    --http.corsdomain '*' \
    --http.addr "0.0.0.0" \
    --nodiscover --maxpeers 0 --mine \
    --networkid $NETWORK_ID \
    --dev \
    --allow-insecure-unlock \
    --rpc.allow-unprotected-txs \
    --dev.gaslimit 12000000)

  sleep 3

  echo -e "Account balances(Defaults):"
  for ACCOUNT in $DEFAULT_ADDRESS_1 $DEFAULT_ADDRESS_2 $DEFAULT_ADDRESS_3; do
    isSigner=" (Default account)"
    if [ "$ACCOUNT" == "$DEFAULT_ADDRESS_1" ]; then
      isSigner=" (Bundler signer account)"
    fi

    docker exec $geth_container_id geth \
      --exec "eth.sendTransaction({from: eth.accounts[0], to: \"$ACCOUNT\", value: web3.toWei(4337, \"ether\")})" \
      attach http://localhost:$GETH_PORT/ > /dev/null
    
    balance=$(docker exec $geth_container_id geth --exec "eth.getBalance(\"$ACCOUNT\")" attach http://localhost:$GETH_PORT/)
    echo -e "  - $ACCOUNT$isSigner: $balance wei"
  done
  echo -e "\n"
}

stop_all() {
  docker stop $geth_container_id > /dev/null
  docker rm $geth_container_id > /dev/null
  exit 0
}

# Start eth-node
start_geth

trap stop_all SIGINT 
echo -e "Press Ctrl+C to stop the eth-node."
while true; do
  sleep 1
done