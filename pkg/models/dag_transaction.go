package models

import "time"

type DagTransaction struct {
	TxHash        string    `json:"txHash" db:"tx_hash"`
	BlockHash     string    `json:"blockHash" db:"block_hash"`
	BlockOrder    uint64    `json:"blockOrder" db:"block_order"`
	TxIndex       uint32    `json:"txIndex" db:"tx_index"`
	Size          int32     `json:"size" db:"size"`
	Confirmations int64     `json:"confirmations" db:"confirmations"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	TxsValid      bool      `json:"txsValid" db:"txs_valid"`
}
