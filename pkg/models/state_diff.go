package models

type StateDiff struct {
	Address string `json:"address" db:"address"`
	Slot    string `json:"slot" db:"slot"`
	Before  string `json:"before" db:"before_value"`
	After   string `json:"after" db:"after_value"`
}
