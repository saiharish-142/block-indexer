package models

type BlockTip struct {
	Hash   string `json:"hash" db:"hash"`
	Height uint64 `json:"height" db:"height"`
	Status string `json:"status" db:"status"` // "valid" or "invalid"
}
