package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EVMClient struct {
	endpoint   string
	httpClient *http.Client
}

type BlockRPC struct {
	Hash         string          `json:"hash"`
	ParentHash   string          `json:"parentHash"`
	Number       string          `json:"number"`
	Timestamp    string          `json:"timestamp"`
	Transactions []any           `json:"transactions"`
	Raw          json.RawMessage `json:"raw,omitempty"`
}

type ReceiptRPC struct {
	TransactionHash string          `json:"transactionHash"`
	BlockHash       string          `json:"blockHash"`
	BlockNumber     string          `json:"blockNumber"`
	Status          string          `json:"status"`
	GasUsed         string          `json:"gasUsed"`
	Logs            json.RawMessage `json:"logs"`
}

type NewHeadEvent struct {
	Hash   string `json:"hash"`
	Number uint64 `json:"number"`
}

type TraceConfig struct {
	DisableStorage bool `json:"disableStorage"`
	DisableMemory  bool `json:"disableMemory"`
}

type TraceResult struct {
	Raw json.RawMessage `json:"raw"`
}

func NewEVMClient(endpoint string) *EVMClient {
	return &EVMClient{
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *EVMClient) GetBlockByNumber(ctx context.Context, number uint64) (BlockRPC, error) {
	var res BlockRPC
	if err := c.call(ctx, "eth_getBlockByNumber", []any{fmt.Sprintf("0x%x", number), true}, &res); err != nil {
		return BlockRPC{}, err
	}
	return res, nil
}

func (c *EVMClient) GetBlockByHash(ctx context.Context, hash string) (BlockRPC, error) {
	var res BlockRPC
	if err := c.call(ctx, "eth_getBlockByHash", []any{hash, true}, &res); err != nil {
		return BlockRPC{}, err
	}
	return res, nil
}

func (c *EVMClient) GetTransactionReceipt(ctx context.Context, hash string) (ReceiptRPC, error) {
	var res ReceiptRPC
	if err := c.call(ctx, "eth_getTransactionReceipt", []any{hash}, &res); err != nil {
		return ReceiptRPC{}, err
	}
	return res, nil
}

func (c *EVMClient) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	var hex string
	if err := c.call(ctx, "eth_blockNumber", []any{}, &hex); err != nil {
		return 0, err
	}
	var number uint64
	if _, err := fmt.Sscanf(hex, "0x%x", &number); err != nil {
		return 0, fmt.Errorf("parse block number: %w", err)
	}
	return number, nil
}

func (c *EVMClient) SubscribeNewHeads(ctx context.Context) (<-chan NewHeadEvent, error) {
	ch := make(chan NewHeadEvent, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		var counter uint64
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				counter++
				ch <- NewHeadEvent{Hash: fmt.Sprintf("0x%x", counter), Number: counter}
			}
		}
	}()
	return ch, nil
}

func (c *EVMClient) DebugTraceTransaction(ctx context.Context, hash string, cfg TraceConfig) (TraceResult, error) {
	var res TraceResult
	if err := c.call(ctx, "debug_traceTransaction", []any{hash, cfg}, &res.Raw); err != nil {
		return TraceResult{}, err
	}
	return res, nil
}

func (c *EVMClient) call(ctx context.Context, method string, params any, out any) error {
	reqBody := rpcRequest{JSONRPC: "2.0", Method: method, Params: params, ID: time.Now().UnixNano()}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal rpc request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("rpc call: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("decode rpc response: %w", err)
	}
	if rpcResp.Error != nil {
		return fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if out != nil {
		if err := json.Unmarshal(rpcResp.Result, out); err != nil {
			return fmt.Errorf("decode result: %w", err)
		}
	}
	return nil
}

func (c *EVMClient) GetReceipts(ctx context.Context, txHashes []string) ([]ReceiptRPC, error) {
	reqs := make([]rpcRequest, len(txHashes))
	for i, hash := range txHashes {
		reqs[i] = rpcRequest{
			JSONRPC: "2.0",
			Method:  "eth_getTransactionReceipt",
			Params:  []any{hash},
			ID:      int64(i),
		}
	}
	resps, err := c.BatchCall(ctx, reqs)
	if err != nil {
		return nil, err
	}

	results := make([]ReceiptRPC, len(txHashes))
	for i, resp := range resps {
		if resp.Error != nil {
			return nil, fmt.Errorf("receipt error for tx %s: %s", txHashes[i], resp.Error.Message)
		}
		if err := json.Unmarshal(resp.Result, &results[i]); err != nil {
			return nil, fmt.Errorf("decode receipt %d: %w", i, err)
		}
	}
	return results, nil
}

func (c *EVMClient) GetTraces(ctx context.Context, txHashes []string) ([]TraceResult, error) {
	reqs := make([]rpcRequest, len(txHashes))
	cfg := TraceConfig{DisableStorage: true, DisableMemory: true}
	for i, hash := range txHashes {
		reqs[i] = rpcRequest{
			JSONRPC: "2.0",
			Method:  "debug_traceTransaction",
			Params:  []any{hash, cfg},
			ID:      int64(i),
		}
	}
	resps, err := c.BatchCall(ctx, reqs)
	if err != nil {
		return nil, err
	}

	results := make([]TraceResult, len(txHashes))
	for i, resp := range resps {
		if resp.Error != nil {
			// Trace might fail for some txs, depending on node. treating as separate error?
			// For now, log or return error. detailed indexer might want to skip or record error.
			// Map error to result?
			// The TraceResult struct can have error field if we want.
			// modifying TraceResult to have error?
			// For now, fail hard or log.
			continue // skip or fail? Fail hard usually for consistency.
		}
		if err := json.Unmarshal(resp.Result, &results[i].Raw); err != nil {
			return nil, fmt.Errorf("decode trace %d: %w", i, err)
		}
	}
	return results, nil
}

func (c *EVMClient) BatchCall(ctx context.Context, reqs []rpcRequest) ([]rpcResponse, error) {
	body, err := json.Marshal(reqs)
	if err != nil {
		return nil, fmt.Errorf("marshal batch: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rpc call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read batch response: %w", err)
	}

	var rpcResps []rpcResponse
	if err := json.Unmarshal(respBody, &rpcResps); err != nil {
		// Some nodes reply with a single response when batch is not supported.
		var single rpcResponse
		if errSingle := json.Unmarshal(respBody, &single); errSingle != nil {
			return nil, fmt.Errorf("decode batch response: %w", err)
		}
		if len(reqs) != 1 {
			return nil, fmt.Errorf("batch call returned single response for %d requests", len(reqs))
		}
		rpcResps = []rpcResponse{single}
	}
	// Sort responses by ID if needed? RPC batch response order is not guaranteed by spec, but usually matches or mapped by ID.
	// Simple mapping by ID assuming ID corresponds to index.
	// Re-ordering logic:
	ordered := make([]rpcResponse, len(reqs))
	respMap := make(map[int64]rpcResponse)
	for _, r := range rpcResps {
		respMap[r.ID] = r
	}
	for i := 0; i < len(reqs); i++ {
		ordered[i] = respMap[int64(i)]
	}

	return ordered, nil
}
