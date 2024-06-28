# Default accounts


## Coinbase Account
In geth dev mode, the term "coinbase" refers to the primary account used for mining rewards and initial transactions. 

A random, pre-allocated developer account will be available and unlocked as `eth.coinbase`, which can be used for testing. You can retrieve the list of accounts and use the first account as the coinbase account. Here’s how you can do it:


1. Get the list of accounts:
```curl
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545
```

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": ["0xYourFirstAccountAddress"]
}
```

2. Get the balance of the first account:
```curl
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xYourFirstAccountAddress", "latest"],"id":1}' -H "Content-Type: application/json" http://localhost:8545
```

The response will look something like this:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x4563918244f40000"
}
```

1. You can run the following script get the coinbase account address and balance:
```bash
make coinbase-account
```