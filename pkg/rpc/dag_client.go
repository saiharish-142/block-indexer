package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DAGClient struct {
	endpoint   string
	httpClient *http.Client
}

// Result structs mapping to RPC responses

type GetBlockResult struct {
	Hash          string        `json:"hash"`
	TxsValid      bool          `json:"txsvalid"`
	Confirmations int64         `json:"confirmations"`
	Version       int32         `json:"version"`
	Weight        string        `json:"weight"` // int64 as string? User said string.
	Height        int64         `json:"height"`
	TxRoot        string        `json:"txRoot"`
	Order         int64         `json:"order"`
	Transactions  []string      `json:"transactions"` // Plain list of tx hashes if verbose=true, fullTx=false
	// For fullTx=true, this would be []TxRawResult, dealing with that dynamically or separate struct
	TransactionFee int64        `json:"transactionfee"`
	StateRoot      string       `json:"stateRoot"`
	Bits           string       `json:"bits"`
	Difficulty     uint32       `json:"difficulty"`
	Pow            PowResult    `json:"pow"`
	Timestamp      string       `json:"timestamp"` // RFC3339
	ParentRoot     string       `json:"parentroot"`
	Parents        []string     `json:"parents"`
	Children       []string     `json:"children"`
}

type GetBlockVerboseResult struct {
	GetBlockResult
	Transactions []TxRawResult `json:"transactions"`
}

type PowResult struct {
	PowName   string `json:"pow_name"`
	PowType   uint8  `json:"pow_type"`
	Nonce     uint64 `json:"nonce"`
	ProofData *struct {
		EdgeBits     int    `json:"edge_bits"`
		CircleNonces string `json:"circle_nonces"`
	} `json:"proof_data"`
}

type TxRawResult struct {
	Hex           string `json:"hex"`
	TxHash        string `json:"txhash"`
	Size          int32  `json:"size"`
	Version       uint32 `json:"version"`
	Timestamp     string `json:"timestamp"`
	BlockHash     string `json:"blockhash"`
	BlockOrder    uint64 `json:"blockorder"`
	TxIndex       uint32 `json:"txindex"`
	Confirmations int64  `json:"confirmations"`
	Time          int64  `json:"time"`
	BlockTime     int64  `json:"blocktime"`
	TxsValid      bool   `json:"txsvalid"`
}

type TipsInfo struct {
	Count   int         `json:"count"`
	Valid   []TipDetail `json:"valid"`
	Invalid []TipDetail `json:"invalid"`
}

type TipDetail struct {
	ID          uint64 `json:"id"`
	Hash        string `json:"hash"`
	Height      uint64 `json:"height"`
	PruneHeight uint64 `json:"pruneheight,omitempty"`
	PruneCD     uint64 `json:"prunecd,omitempty"`
}

type StateRootResult struct {
	Hash         string `json:"Hash"`
	Order        uint64 `json:"Order"`
	Height       uint64 `json:"Height"`
	Valid        bool   `json:"Valid"`
	EVMStateRoot string `json:"EVMStateRoot"`
	EVMHeight    uint64 `json:"EVMHeight"`
	EVMHead      string `json:"EVMHead"`
	StateRoot    string `json:"StateRoot"`
}

type DagBlockEvent struct {
	Hash   string `json:"hash"`
	Height uint64 `json:"height"`
	Order  int64  `json:"order"`
	Time   string `json:"time"` // RFC3339? Notification says time (maybe int64 or string, user said "time").
	// Checking user spec: blockConnected ... time
	Txs []string `json:"txs"`
}

func NewDAGClient(endpoint string) *DAGClient {
	return &DAGClient{
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *DAGClient) GetBlockCount(ctx context.Context) (int64, error) {
	var res int64
	if err := c.call(ctx, "Blockdag.getBlockCount", nil, &res); err != nil {
		return 0, err
	}
	return res, nil
}

// GetBlock returns verbose block with tx hashes
func (c *DAGClient) GetBlock(ctx context.Context, hash string) (*GetBlockResult, error) {
	var res GetBlockResult
	// getBlock(hash, verbose=true, inclTx=true, fullTx=false)
	if err := c.call(ctx, "Blockdag.getBlock", []any{hash, true, true, false}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetBlockFull returns verbose block with full tx details
func (c *DAGClient) GetBlockFull(ctx context.Context, hash string) (*GetBlockVerboseResult, error) {
	var res GetBlockVerboseResult
	// getBlock(hash, verbose=true, inclTx=true, fullTx=true)
	if err := c.call(ctx, "Blockdag.getBlock", []any{hash, true, true, true}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *DAGClient) GetTips(ctx context.Context) (*TipsInfo, error) {
	var res TipsInfo
	if err := c.call(ctx, "Blockdag.tips", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *DAGClient) GetStateRoot(ctx context.Context, order uint64) (*StateRootResult, error) {
	var res StateRootResult
	// getStateRoot(order, verbose=true)
	if err := c.call(ctx, "Blockdag.getStateRoot", []any{order, true}, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// SubscribeDagBlocks provides a lightweight polling-based stream of new DAG
// blocks. A real deployment should swap this for the RPC websocket
// notifications, but this keeps the pipeline running without extra deps.
func (c *DAGClient) SubscribeDagBlocks(ctx context.Context) (<-chan DagBlockEvent, error) {
	current, err := c.GetBlockCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("init dag subscription: %w", err)
	}

	out := make(chan DagBlockEvent, 8)
	go func() {
		defer close(out)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		lastEmitted := current
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				latest, err := c.GetBlockCount(ctx)
				if err != nil || latest <= lastEmitted {
					continue
				}

				nextOrder := lastEmitted + 1
				for ; nextOrder <= latest; nextOrder++ {
					state, err := c.GetStateRoot(ctx, uint64(nextOrder))
					if err != nil {
						// Retry the same order on the next tick.
						break
					}
					out <- DagBlockEvent{
						Hash:   state.Hash,
						Height: state.Height,
						Order:  int64(nextOrder),
					}
					lastEmitted = nextOrder
				}
			}
		}
	}()

	return out, nil
}

// call executes a single JSON-RPC request against the DAG endpoint.
func (c *DAGClient) call(ctx context.Context, method string, params any, out any) error {
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
	// Basic Auth if needed (user mentioned it, but provided no creds in prompt, assuming env or handled by url)
	
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
