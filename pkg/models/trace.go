package models

type TxTrace struct {
	TxHash string `json:"txHash" db:"tx_hash"`
	Trace  any    `json:"trace" db:"trace"`
}
