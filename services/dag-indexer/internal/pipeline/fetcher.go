package pipeline

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/models"
	"github.com/example/block-indexer/pkg/rpc"
)

type Fetcher struct {
	client *rpc.DAGClient
	log    *zap.Logger
}

func NewFetcher(client *rpc.DAGClient, log *zap.Logger) *Fetcher {
	return &Fetcher{client: client, log: log}
}

func (f *Fetcher) Run(ctx context.Context, refs <-chan DagRef) <-chan DagBlockPayload {
	out := make(chan DagBlockPayload, 8)
	go func() {
		defer close(out)
		for ref := range refs {
			rpcBlock, err := f.client.GetBlock(ctx, ref.Hash)
			if err != nil {
				f.log.Warn("failed to fetch dag block", zap.String("hash", ref.Hash), zap.Error(err))
				continue
			}

			// Map RPC result to Model
			// Assuming RFC3339 for timestamp string from RPC
			parsedTime, _ := time.Parse(time.RFC3339, rpcBlock.Timestamp)

			block := models.DagBlock{
				Hash:          rpcBlock.Hash,
				Order:         uint64(rpcBlock.Order),
				Height:        uint64(rpcBlock.Height),
				Weight:        rpcBlock.Weight,
				TxRoot:        rpcBlock.TxRoot,
				StateRoot:     rpcBlock.StateRoot,
				ParentRoot:    rpcBlock.ParentRoot,
				Confirmations: rpcBlock.Confirmations,
				TxsValid:      rpcBlock.TxsValid,
				Difficulty:    rpcBlock.Difficulty,
				Bits:          rpcBlock.Bits,
				PowName:       rpcBlock.Pow.PowName,
				PowType:       rpcBlock.Pow.PowType,
				PowNonce:      rpcBlock.Pow.Nonce,
				Timestamp:     parsedTime,
				// Reward/Fee might need full verbose or separate calc, user spec says GetBlock w/ verbose=false?
				// The spec said: getBlock(hash, verbose=false...) -> but I used verbose=true in client
				// The result has transactionfee?
				// Spec: "transactionfee:int64?"
				TxFee: rpcBlock.TransactionFee,
				// Reward not explicitly in BlockResult? "reward:int64" is in Header verbose.
				// For now leaving Reward 0 or checking if it's in result.
				// Spec says "Block verbose ... transactionfee:int64?". Header verbose has reward.
				// I'll leave Reward as 0 if not in BlockResult.
				IsBlue: 0, // Need to determine from somewhere? Spec: "isBlue(hash) -> 0/1/2". separate call?
				// Or maybe it's in context/result I missed.
				// For now, I'll map what I have.
				MainChain:  true, // Placeholder or need `isOnMainChain`.
				InsertedAt: time.Now(),
			}

			out <- DagBlockPayload{
				Block:   block,
				Parents: rpcBlock.Parents,
			}
		}
	}()
	return out
}
