# Mempool

Besty maintains in memory mempool by pooling the bundlers `debug_bundler_dumpMempool` rpc method.

Example:
```shell
curl -X POST --data '{"jsonrpc":"2.0","method":"debug_bundler_dumpMempool","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:4337/rpc
```

result:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": [
    {
      "sender": "0x0B09809bE0Cc9FA0938C934F5845A98fFfcAA155",
      "nonce": "0x00",
      "factory": "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
      "factoryData": "0x5fbfb9cf00000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c8000000000000000000000000000000000000000000000000000000000016d826",
      "callData": "0xb61d27f60000000000000000000000009fe46736679d2d9a65f0992f2272de9f3c7fa6e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000004d09de08a00000000000000000000000000000000000000000000000000000000",
      "signature": "0x124ac5d541198d44fc334d7aba6fb982cd2be669fc2b7b071623f1855bf4daf34bb22a81e19cd030b6fcade74447e863fddaaac6da749cffa46aea61ad3fc6671b",
      "callGasLimit": "0x574e",
      "verificationGasLimit": "0x0f4240",
      "preVerificationGas": "0xaf30",
      "maxPriorityFeePerGas": "0x59682f00",
      "maxFeePerGas": "0x667b4fca"
    }
  ]
}
```