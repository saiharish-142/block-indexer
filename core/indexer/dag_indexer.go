package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/example/block-indexer/core/pb"
	"go.uber.org/zap"
)

// fetchDagBlockByOrder calls a DAG-style RPC with basic auth to fetch a block by order.
func (i *Indexer) fetchDagBlockByOrder(ctx context.Context, order uint64, verbose, inclTx, fullTx bool) (*pb.BlockSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	i.logger.Info("Order", zap.Uint64("order", order))

	params := []any{
		order,
		verbose,
		inclTx,
		fullTx,
	}

	reqBody, err := json.Marshal(dagRPCRequest{
		JSONRPC: "2.0",
		Method:  "getBlockByOrder",
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal dag rpc request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.cfg.DagRPCURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("build dag rpc request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(i.cfg.DagRPCUser, i.cfg.DagRPCPass)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call dag rpc: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// i.logger.Info("DAG RPC Response", zap.ByteString("body", body))

	var rpcResp dagRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("decode dag rpc response: %w", err)
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("dag rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if rpcResp.Result == nil {
		return nil, errors.New("dag rpc returned no result")
	}

	orderNum, err := parseUintFromAny(rpcResp.Result["order"])
	if err != nil {
		return nil, fmt.Errorf("parse dag order: %w", err)
	}
	ts, _ := parseTimestamp(rpcResp.Result["timestamp"])

	hash, _ := rpcResp.Result["hash"].(string)
	parent, _ := rpcResp.Result["parentHash"].(string)
	if parent == "" {
		parent, _ = rpcResp.Result["previousHash"].(string)
	}
	if parent == "" {
		parent, _ = rpcResp.Result["parentroot"].(string)
	}
	if parent == "" {
		if parents, ok := rpcResp.Result["parents"].([]any); ok && len(parents) > 0 {
			if p, ok := parents[0].(string); ok && p != "null" {
				parent = p
			}
		}
	}

	return &pb.BlockSummary{
		Number:     orderNum,
		Hash:       hash,
		ParentHash: parent,
		Timestamp:  int64(ts),
	}, nil
}

func parseUintFromAny(v any) (uint64, error) {
	switch val := v.(type) {
	case nil:
		return 0, fmt.Errorf("value is nil")
	case string:
		s := strings.TrimPrefix(val, "0x")
		if strings.HasPrefix(val, "0x") {
			return strconv.ParseUint(s, 16, 64)
		}
		return strconv.ParseUint(val, 10, 64)
	case float64:
		return uint64(val), nil
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}

func parseTimestamp(v any) (int64, error) {
	switch val := v.(type) {
	case string:
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t.Unix(), nil
		}
		if n, err := parseUintFromAny(val); err == nil {
			return int64(n), nil
		}
		return 0, fmt.Errorf("unsupported timestamp string %q", val)
	case float64:
		return int64(val), nil
	default:
		if n, err := parseUintFromAny(val); err == nil {
			return int64(n), nil
		}
		return 0, fmt.Errorf("unsupported timestamp type %T", v)
	}
}

type dagRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type dagRPCResponse struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      int            `json:"id"`
	Result  map[string]any `json:"result"`
	Error   *dagRPCError   `json:"error"`
}

type dagRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
