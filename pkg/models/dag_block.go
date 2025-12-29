package models

import "time"

type DagBlock struct {
	Hash          string    `json:"hash" db:"hash"`
	Order         uint64    `json:"order" db:"block_order"`
	Height        uint64    `json:"height" db:"height"`
	Weight        string    `json:"weight" db:"weight"`
	TxRoot        string    `json:"txRoot" db:"tx_root"`
	StateRoot     string    `json:"stateRoot" db:"state_root"`
	ParentRoot    string    `json:"parentRoot" db:"parent_root"`
	Confirmations int64     `json:"confirmations" db:"confirmations"`
	TxsValid      bool      `json:"txsValid" db:"txs_valid"`
	Difficulty    uint32    `json:"difficulty" db:"difficulty"`
	Bits          string    `json:"bits" db:"bits"`
	PowName       string    `json:"powName" db:"pow_name"`
	PowType       uint8     `json:"powType" db:"pow_type"`
	PowNonce      uint64    `json:"powNonce" db:"pow_nonce"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	Reward        int64     `json:"reward" db:"reward"`
	TxFee         int64     `json:"txFee" db:"tx_fee"`
	IsBlue        uint8     `json:"isBlue" db:"is_blue"`
	MainChain     bool      `json:"mainChain" db:"main_chain"`
	InsertedAt    time.Time `json:"insertedAt" db:"inserted_at"`
}
