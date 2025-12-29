package models

import "time"

type Transaction struct {
	Hash            string    `json:"hash" db:"hash"`
	BlockHash       string    `json:"blockHash" db:"block_hash"`
	BlockNumber     uint64    `json:"blockNumber" db:"block_number"`
	Index           int       `json:"index" db:"tx_index"`
	From            string    `json:"from" db:"from_address"`
	To              string    `json:"to,omitempty" db:"to_address"`
	Nonce           uint64    `json:"nonce" db:"nonce"`
	Value           string    `json:"value" db:"value"`
	GasPrice        string    `json:"gasPrice" db:"gas_price"`
	GasUsed         string    `json:"gasUsed" db:"gas_used"`
	GasLimit        string    `json:"gasLimit" db:"gas_limit"`
	Status          int       `json:"status" db:"status"`
	MethodSignature string    `json:"methodSignature,omitempty" db:"method_sig"`
	Input           string    `json:"input,omitempty" db:"input"`
	ContractAddress string    `json:"contractAddress,omitempty" db:"contract_address"`
	Type            int       `json:"type" db:"type"`
	IsCanonical     bool      `json:"isCanonical" db:"is_canonical"`
	Timestamp       time.Time `json:"timestamp" db:"timestamp"`
}
