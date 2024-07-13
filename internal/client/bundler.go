package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/transeptorlabs/betsy/internal/data"
)

// BundlerClient is a client for the bundler node
type BundlerClient struct {
	bundlerUrl       string
	jsonRpcRequestID int
	mutex            sync.Mutex
}

// jsonrpcBase is the base struct for json rpc requests
type jsonrpcBase struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
}

// debugBundlerDumpMempoolRes is the response struct for debug_bundler_dumpMempool rpc method
type debugBundlerDumpMempoolRes struct {
	jsonrpcBase
	Result []data.UserOpV7Hexify `json:"result"`
}

// debug_bundler_addUserOpsRes is the response struct for debug_bundler_addUserOps rpc method
type debug_bundler_addUserOpsRes struct {
	jsonrpcBase
	Result string `json:"result"`
}

// NewBundlerClient creates a new BundlerClient
func NewBundlerClient(bundlerUrl string) *BundlerClient {
	return &BundlerClient{
		bundlerUrl:       bundlerUrl,
		jsonRpcRequestID: 1,
	}
}

// getRequest creates a new http request for the given rpc method and params
func (b *BundlerClient) getRequest(rpcMethod string, params []interface{}) (*http.Request, error) {
	// Make json rpc request
	jsonBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      b.jsonRpcRequestID,
		"method":  rpcMethod,
		"params":  params,
	})
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, b.bundlerUrl, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Debug_bundler_dumpMempool calls the debug_bundler_dumpMempool rpc method
func (b *BundlerClient) Debug_bundler_dumpMempool() ([]data.UserOpV7Hexify, error) {
	log.Debug().Msgf("Making call to bundler node debug_bundler_dumpMempool at %s", b.bundlerUrl)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	req, err := b.getRequest("debug_bundler_dumpMempool", []interface{}{})
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// handle json rpc response
	b.jsonRpcRequestID = b.jsonRpcRequestID + 1
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Request to bundler debug_bundler_dumpMempool rpc method failed with status code: %d", res.StatusCode))
	}

	resJsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// parse result
	var data *debugBundlerDumpMempoolRes
	err = json.Unmarshal(resJsonBody, &data)
	if err != nil {
		return nil, err
	}

	return data.Result, nil
}

// Debug_bundler_addUserOps calls the debug_bundler_addUserOps rpc method
func (b *BundlerClient) Debug_bundler_addUserOps(ops []data.UserOpV7Hexify) error {
	log.Debug().Msgf("Making call to bundler node debug_bundler_addUsers at %s", b.bundlerUrl)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if ops == nil || len(ops) == 0 {
		return errors.New("Can not add empty userOps")
	}

	req, err := b.getRequest("debug_bundler_addUserOps", []interface{}{
		ops,
	})
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// handle json rpc response
	b.jsonRpcRequestID = b.jsonRpcRequestID + 1
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Request to bundler debug_bundler_addUserOps rpc method failed with status code: %d", res.StatusCode))
	}

	resJsonBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// parse result
	var data *debug_bundler_addUserOpsRes
	err = json.Unmarshal(resJsonBody, &data)
	if err != nil {
		return err
	}

	return nil
}
