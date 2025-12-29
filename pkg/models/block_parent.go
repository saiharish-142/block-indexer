package models

type BlockParent struct {
	BlockHash  string `json:"blockHash" db:"block_hash"`
	ParentHash string `json:"parentHash" db:"parent_hash"`
}
