package models

import "time"

type InternalTransaction struct {
	TxHash       string    `json:"txHash" db:"tx_hash"`
	BlockHash    string    `json:"blockHash" db:"block_hash"`
	BlockNumber  uint64    `json:"blockNumber" db:"block_number"`
	From         string    `json:"from" db:"from_address"`
	To           string    `json:"to" db:"to_address"`
	Value        string    `json:"value" db:"value"`
	GasLimit     uint64    `json:"gasLimit" db:"gas_limit"`
	GasUsed      uint64    `json:"gasUsed" db:"gas_used"`
	CallType     string    `json:"callType" db:"call_type"`
	TraceAddress string    `json:"traceAddress" db:"trace_address"`
	Error        string    `json:"error,omitempty" db:"error"`
	IsCanonical  bool      `json:"isCanonical" db:"is_canonical"`
	InsertedAt   time.Time `json:"insertedAt" db:"inserted_at"`
}
