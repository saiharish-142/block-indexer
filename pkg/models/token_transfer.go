package models

import "time"

type TokenTransfer struct {
	TxHash       string    `json:"txHash" db:"tx_hash"`
	BlockHash    string    `json:"blockHash" db:"block_hash"`
	BlockNumber  uint64    `json:"blockNumber" db:"block_number"`
	LogIndex     int       `json:"logIndex" db:"log_index"`
	TokenAddress string    `json:"tokenAddress" db:"token_address"`
	From         string    `json:"from" db:"from_address"`
	To           string    `json:"to" db:"to_address"`
	Value        string    `json:"value" db:"value"`
	TokenID      string    `json:"tokenId" db:"token_id"`
	Type         string    `json:"type" db:"type"` // ERC20, ERC721
	IsCanonical  bool      `json:"isCanonical" db:"is_canonical"`
	InsertedAt   time.Time `json:"insertedAt" db:"inserted_at"`
}
