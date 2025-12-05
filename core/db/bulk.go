package db

import (
    "context"
    "time"

    "github.com/example/block-indexer/core/pb"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// CopyBlocks ingests a slice of blocks into Postgres using CopyFrom for throughput.
func CopyBlocks(ctx context.Context, pool *pgxpool.Pool, blocks []pb.BlockSummary) error {
    rows := make([][]any, 0, len(blocks))
    for _, b := range blocks {
        rows = append(rows, []any{
            b.Number,
            b.Hash,
            b.ParentHash,
            time.Unix(b.Timestamp, 0).UTC(),
            int64(b.GasUsed),
            int64(b.GasLimit),
            b.Miner,
            b.Nonce,
            b.Difficulty,
            b.ExtraData,
            b.LogsBloom,
            b.MixHash,
            b.ReceiptsRoot,
            b.Sha3Uncles,
            int64(b.SizeBytes),
            b.StateRoot,
            b.TxRoot,
            int64(b.TxCount),
            b.Uncles,
            b.TxHashes,
        })
    }

    _, err := pool.CopyFrom(
        ctx,
        pgx.Identifier{"blocks"},
        []string{
            "number",
            "hash",
            "parent_hash",
            "timestamp",
            "gas_used",
            "gas_limit",
            "miner",
            "nonce",
            "difficulty",
            "extra_data",
            "logs_bloom",
            "mix_hash",
            "receipts_root",
            "sha3_uncles",
            "size_bytes",
            "state_root",
            "tx_root",
            "tx_count",
            "uncles",
            "tx_hashes",
        },
        pgx.CopyFromRows(rows),
    )
    return err
}

// CopyDagBlocks ingests DAG blocks into Postgres using CopyFrom.
func CopyDagBlocks(ctx context.Context, pool *pgxpool.Pool, blocks []pb.BlockSummary) error {
    rows := make([][]any, 0, len(blocks))
    for _, b := range blocks {
        rows = append(rows, []any{b.Number, b.Hash, b.ParentHash, time.Unix(b.Timestamp, 0).UTC()})
    }

    _, err := pool.CopyFrom(
        ctx,
        pgx.Identifier{"dag_blocks"},
        []string{"number", "hash", "parent_hash", "timestamp"},
        pgx.CopyFromRows(rows),
    )
    return err
}
