package models

import "time"

type StatsOverview struct {
	BlockCount       uint64 `json:"blockCount"`
	TransactionCount uint64 `json:"transactionCount"`
	AvgGasPrice      string `json:"avgGasPrice"`
}

type HistoricalPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
	Metric    string    `json:"metric"`
}
