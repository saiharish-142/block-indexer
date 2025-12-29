package workers

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/models"
	"github.com/example/block-indexer/pkg/rpc"
	"github.com/example/block-indexer/services/trace-service/internal"
	"github.com/example/block-indexer/services/trace-service/internal/repository"
)

type TraceWorker struct {
	cfg      internal.Config
	log      *zap.Logger
	repo     *repository.TraceRepo
	diffRepo *repository.StateDiffRepo
	client   *rpc.EVMClient
}

func NewTraceWorker(cfg internal.Config, log *zap.Logger, repo *repository.TraceRepo, diffRepo *repository.StateDiffRepo, client *rpc.EVMClient) *TraceWorker {
	return &TraceWorker{
		cfg:      cfg,
		log:      log,
		repo:     repo,
		diffRepo: diffRepo,
		client:   client,
	}
}

func (w *TraceWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.drain(ctx)
			}
		}
	}()
}

func (w *TraceWorker) drain(ctx context.Context) {
	hashes, err := w.repo.NextQueued(ctx, w.cfg.Trace.Concurrency)
	if err != nil {
		w.log.Warn("failed to fetch queued traces", zap.Error(err))
		return
	}
	for _, hash := range hashes {
		go w.process(ctx, hash)
	}
}

func (w *TraceWorker) process(ctx context.Context, hash string) {
	if err := w.repo.MarkProcessing(ctx, hash); err != nil {
		w.log.Warn("failed to mark processing", zap.String("hash", hash), zap.Error(err))
		return
	}
	trace, err := w.client.DebugTraceTransaction(ctx, hash, rpc.TraceConfig{})
	if err != nil {
		w.log.Warn("trace failed", zap.String("hash", hash), zap.Error(err))
		return
	}
	if err := w.repo.StoreTrace(ctx, models.TxTrace{TxHash: hash, Trace: trace.Raw}); err != nil {
		w.log.Warn("failed to store trace", zap.String("hash", hash), zap.Error(err))
	}
}
