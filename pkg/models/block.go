package models

import "time"

type Block struct {
	Hash        string    `json:"hash" db:"hash"`
	ParentHash  string    `json:"parentHash" db:"parent_hash"`
	Number      uint64    `json:"number" db:"number"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Miner       string    `json:"miner" db:"miner"`
	GasUsed     string    `json:"gasUsed" db:"gas_used"`
	GasLimit    string    `json:"gasLimit" db:"gas_limit"`
	TxCount     int       `json:"txCount" db:"tx_count"`
	IsCanonical bool      `json:"isCanonical" db:"is_canonical"`
	Raw         any       `json:"raw" db:"raw"`
}

type BlockEnvelope struct {
	Block          Block                 `json:"block"`
	Transactions   []Transaction         `json:"transactions"`
	Logs           []Log                 `json:"logs"`
	InternalTxs    []InternalTransaction `json:"internalTransactions"`
	TokenTransfers []TokenTransfer       `json:"tokenTransfers"`
}
