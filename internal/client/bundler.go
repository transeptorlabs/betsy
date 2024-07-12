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

type BundlerClient struct {
	bundlerUrl       string
	jsonRpcRequestID int
	mutex            sync.Mutex
}

type jsonrpcBase struct {
	jsonrpc string
	id      int
}

type debugBundlerDumpMempoolRes struct {
	jsonrpcBase
	result []data.UserOpV7Hexify
}

type debug_bundler_addUserOpsRes struct {
	jsonrpcBase
	result string
}

func NewBundlerClient(bundlerUrl string) *BundlerClient {
	return &BundlerClient{
		bundlerUrl:       bundlerUrl,
		jsonRpcRequestID: 1,
	}
}

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

func (b *BundlerClient) Debug_bundler_dumpMempool() ([]data.UserOpV7Hexify, error) {
	log.Info().Msgf("Making call to bundler node debug_bundler_dumpMempool at %s", b.bundlerUrl)
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

	return data.result, nil
}

func (b *BundlerClient) Debug_bundler_addUserOps(ops []data.UserOpV7Hexify) error {
	log.Info().Msgf("Making call to bundler node debug_bundler_addUsers at %s", b.bundlerUrl)
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
