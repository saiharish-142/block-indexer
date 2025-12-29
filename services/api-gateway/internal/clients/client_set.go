package clients

import "github.com/jackc/pgx/v5/pgxpool"

type ClientSet struct {
	EVM      *EVMIndexerClient
	DAG      *DAGIndexerClient
	Stats    *StatsClient
	Search   *SearchClient
	Contract *ContractClient
	Trace    *TraceClient
}

func NewClientSet(pool *pgxpool.Pool) *ClientSet {
	return &ClientSet{
		EVM:      &EVMIndexerClient{pool: pool},
		DAG:      &DAGIndexerClient{pool: pool},
		Stats:    &StatsClient{pool: pool},
		Search:   &SearchClient{pool: pool},
		Contract: &ContractClient{pool: pool},
		Trace:    &TraceClient{pool: pool},
	}
}
