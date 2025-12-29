package models

import "time"

type Log struct {
	TxHash      string    `json:"txHash" db:"tx_hash"`
	BlockHash   string    `json:"blockHash" db:"block_hash"`
	BlockNumber uint64    `json:"blockNumber" db:"block_number"`
	TxIndex     uint      `json:"txIndex" db:"tx_index"`
	LogIndex    uint      `json:"logIndex" db:"log_index"`
	Index       uint      `json:"index" db:"log_index"`
	Address     string    `json:"address" db:"address"`
	Topics      []string  `json:"topics" db:"topics"`
	Data        string    `json:"data" db:"data"`
	IsCanonical bool      `json:"isCanonical" db:"is_canonical"`
	InsertedAt  time.Time `json:"insertedAt" db:"inserted_at"`
}
