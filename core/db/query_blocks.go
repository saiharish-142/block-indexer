package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/example/block-indexer/core/pb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const evmBlockColumns = `number, hash, parent_hash, timestamp, gas_used, gas_limit, miner, nonce, difficulty, extra_data, logs_bloom, mix_hash, receipts_root, sha3_uncles, size_bytes, state_root, tx_root, tx_count, uncles, tx_hashes`

// ListEVMBlocks returns EVM blocks in descending order with simple cursor pagination.
func ListEVMBlocks(ctx context.Context, pool *pgxpool.Pool, limit int, before *uint64) ([]pb.BlockSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pageSize := limit + 1 // fetch one extra to signal a next cursor
	if pageSize < 1 {
		pageSize = 1
	}

	var (
		rows pgx.Rows
		err  error
	)

	if before != nil {
		rows, err = pool.Query(ctx,
			fmt.Sprintf(`SELECT %s FROM blocks WHERE number < $1 ORDER BY number DESC LIMIT $2`, evmBlockColumns),
			*before, pageSize)
	} else {
		rows, err = pool.Query(ctx,
			fmt.Sprintf(`SELECT %s FROM blocks ORDER BY number DESC LIMIT $1`, evmBlockColumns),
			pageSize)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blocks := make([]pb.BlockSummary, 0, pageSize)
	for rows.Next() {
		var (
			number       int64
			hash         string
			parent       string
			ts           time.Time
			gasUsed      sql.NullInt64
			gasLimit     sql.NullInt64
			miner        sql.NullString
			nonce        sql.NullString
			difficulty   sql.NullString
			extraData    sql.NullString
			logsBloom    sql.NullString
			mixHash      sql.NullString
			receiptsRoot sql.NullString
			sha3Uncles   sql.NullString
			sizeBytes    sql.NullInt64
			stateRoot    sql.NullString
			txRoot       sql.NullString
			txCount      sql.NullInt64
			uncles       []string
			txHashes     []string
		)
		if err := rows.Scan(
			&number,
			&hash,
			&parent,
			&ts,
			&gasUsed,
			&gasLimit,
			&miner,
			&nonce,
			&difficulty,
			&extraData,
			&logsBloom,
			&mixHash,
			&receiptsRoot,
			&sha3Uncles,
			&sizeBytes,
			&stateRoot,
			&txRoot,
			&txCount,
			&uncles,
			&txHashes,
		); err != nil {
			return nil, err
		}

		blocks = append(blocks, pb.BlockSummary{
			Number:       uint64(number),
			Hash:         hash,
			ParentHash:   parent,
			Timestamp:    ts.Unix(),
			GasUsed:      asUint64(gasUsed),
			GasLimit:     asUint64(gasLimit),
			Miner:        miner.String,
			Nonce:        nonce.String,
			Difficulty:   difficulty.String,
			ExtraData:    extraData.String,
			LogsBloom:    logsBloom.String,
			MixHash:      mixHash.String,
			ReceiptsRoot: receiptsRoot.String,
			Sha3Uncles:   sha3Uncles.String,
			SizeBytes:    asUint64(sizeBytes),
			StateRoot:    stateRoot.String,
			TxRoot:       txRoot.String,
			TxCount:      asInt(txCount),
			Uncles:       uncles,
			TxHashes:     txHashes,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

// ListDagBlocks returns DAG blocks in descending order with simple cursor pagination.
func ListDagBlocks(ctx context.Context, pool *pgxpool.Pool, limit int, before *uint64) ([]pb.BlockSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pageSize := limit + 1 // fetch one extra to signal a next cursor
	if pageSize < 1 {
		pageSize = 1
	}

	var (
		rows pgx.Rows
		err  error
	)

	if before != nil {
		rows, err = pool.Query(ctx,
			`SELECT number, hash, parent_hash, timestamp FROM dag_blocks WHERE number < $1 ORDER BY number DESC LIMIT $2`,
			*before, pageSize)
	} else {
		rows, err = pool.Query(ctx,
			`SELECT number, hash, parent_hash, timestamp FROM dag_blocks ORDER BY number DESC LIMIT $1`,
			pageSize)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blocks := make([]pb.BlockSummary, 0, pageSize)
	for rows.Next() {
		var (
			number int64
			hash   string
			parent string
			ts     time.Time
		)
		if err := rows.Scan(&number, &hash, &parent, &ts); err != nil {
			return nil, err
		}

		blocks = append(blocks, pb.BlockSummary{
			Number:     uint64(number),
			Hash:       hash,
			ParentHash: parent,
			Timestamp:  ts.Unix(),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func asUint64(v sql.NullInt64) uint64 {
	if !v.Valid || v.Int64 < 0 {
		return 0
	}
	return uint64(v.Int64)
}

func asInt(v sql.NullInt64) int {
	if !v.Valid {
		return 0
	}
	return int(v.Int64)
}
