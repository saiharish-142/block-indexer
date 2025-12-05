package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/example/block-indexer/core/pb"
	"go.uber.org/zap"
)

// fetchEthBlockByNumber fetches a specific block by number over HTTP RPC.
func (i *Indexer) fetchEthBlockByNumber(ctx context.Context, number uint64) (*pb.BlockSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	reqBody, err := json.Marshal(ethRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{fmt.Sprintf("0x%x", number), false},
		ID:      1,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal rpc request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.cfg.ChainRPCURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call rpc: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp ethRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("decode rpc response: %w", err)
	}
	if rpcResp.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}
	if rpcResp.Result == nil {
		return nil, errors.New("rpc returned no result")
	}

	num, err := parseHexUint64(rpcResp.Result.Number)
	if err != nil {
		return nil, fmt.Errorf("parse block number: %w", err)
	}
	ts, err := parseHexUint64(rpcResp.Result.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("parse block timestamp: %w", err)
	}

	gasUsed := parseHexUint64Default(rpcResp.Result.GasUsed)
	gasLimit := parseHexUint64Default(rpcResp.Result.GasLimit)
	sizeBytes := parseHexUint64Default(rpcResp.Result.Size)
	txHashes := extractTxHashes(rpcResp.Result.Transactions)

	return &pb.BlockSummary{
		Number:       num,
		Hash:         rpcResp.Result.Hash,
		Miner:        rpcResp.Result.Miner,
		ParentHash:   rpcResp.Result.ParentHash,
		Timestamp:    int64(ts),
		GasUsed:      gasUsed,
		GasLimit:     gasLimit,
		Nonce:        rpcResp.Result.Nonce,
		Difficulty:   rpcResp.Result.Difficulty,
		ExtraData:    rpcResp.Result.ExtraData,
		LogsBloom:    rpcResp.Result.LogsBloom,
		MixHash:      rpcResp.Result.MixHash,
		ReceiptsRoot: rpcResp.Result.ReceiptsRoot,
		Sha3Uncles:   rpcResp.Result.Sha3Uncles,
		SizeBytes:    sizeBytes,
		StateRoot:    rpcResp.Result.StateRoot,
		TxRoot:       rpcResp.Result.TransactionsRoot,
		TxCount:      len(txHashes),
		Uncles:       rpcResp.Result.Uncles,
		TxHashes:     txHashes,
	}, nil
}

func parseHexUint64(hexStr string) (uint64, error) {
	trimmed := strings.TrimPrefix(hexStr, "0x")
	return strconv.ParseUint(trimmed, 16, 64)
}

func parseHexUint64Default(hexStr string) uint64 {
	val, err := parseHexUint64(hexStr)
	if err != nil {
		return 0
	}
	return val
}

func extractTxHashes(txField []any) []string {
	if len(txField) == 0 {
		return nil
	}
	hashes := make([]string, 0, len(txField))
	for _, tx := range txField {
		switch v := tx.(type) {
		case string:
			hashes = append(hashes, v)
		case map[string]any:
			if h, ok := v["hash"].(string); ok {
				hashes = append(hashes, h)
			}
		}
	}
	return hashes
}

type ethRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type ethRPCResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      int          `json:"id"`
	Result  *ethRPCBlock `json:"result"`
	Error   *ethRPCError `json:"error"`
}

type ethRPCBlock struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	Miner            string   `json:"miner"`
	ParentHash       string   `json:"parentHash"`
	Timestamp        string   `json:"timestamp"`
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	LogsBloom        string   `json:"logsBloom"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
	Transactions     []any    `json:"transactions"`
}

type ethRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// streamEthHeads subscribes to newHeads over WebSocket and logs them.
func (i *Indexer) streamEthHeads(ctx context.Context) {
	if i.cfg.ChainWSURL == "" {
		return
	}

	conn, _, err := websocket.Dial(ctx, i.cfg.ChainWSURL, nil)
	if err != nil {
		i.logger.Warn("eth ws dial failed", zap.Error(err))
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "shutdown")

	subMsg, _ := json.Marshal(ethWSRequest{
		JSONRPC: "2.0",
		Method:  "eth_subscribe",
		Params:  []any{"newHeads"},
		ID:      1,
	})

	if err := conn.Write(ctx, websocket.MessageText, subMsg); err != nil {
		i.logger.Warn("eth ws subscribe failed", zap.Error(err))
		return
	}

	for {
		var data []byte
		if _, data, err = conn.Read(ctx); err != nil {
			if errors.Is(err, context.Canceled) || websocket.CloseStatus(err) != -1 {
				return
			}
			i.logger.Warn("eth ws read failed", zap.Error(err))
			return
		}

		var msg ethWSMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			i.logger.Warn("eth ws decode failed", zap.Error(err))
			continue
		}

		if msg.Params == nil || msg.Params.Result == nil {
			continue
		}
		head := msg.Params.Result
		num, _ := parseHexUint64(head.Number)
		i.logger.Info("eth head",
			zap.Uint64("number", num),
			zap.String("hash", head.Hash),
			zap.String("parent", head.ParentHash),
		)
	}
}

type ethWSRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type ethWSMessage struct {
	JSONRPC string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  *ethWSParams `json:"params"`
}

type ethWSParams struct {
	Subscription string       `json:"subscription"`
	Result       *ethRPCBlock `json:"result"`
}
