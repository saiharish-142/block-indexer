package models

type Contract struct {
	Address         string `json:"address" db:"address"`
	CreatorTxHash   string `json:"creatorTxHash" db:"creator_tx_hash"`
	CreatedInBlock  uint64 `json:"createdInBlock" db:"created_in_block"`
	Name            string `json:"name" db:"name"`
	Symbol          string `json:"symbol" db:"symbol"`
	Decimals        int    `json:"decimals" db:"decimals"`
	Verified        bool   `json:"verified" db:"verified"`
	ABI             any    `json:"abi" db:"abi"`
	Bytecode        string `json:"bytecode" db:"bytecode"`
	RuntimeBytecode string `json:"runtimeBytecode" db:"runtime_bytecode"`
}
