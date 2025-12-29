package pipeline

import (
	"context"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/events"
	"github.com/example/block-indexer/pkg/models"
	"github.com/example/block-indexer/pkg/rpc"
	"github.com/example/block-indexer/services/evm-indexer/internal"
	"github.com/example/block-indexer/services/evm-indexer/internal/repository"
)

type Manager struct {
	cfg        internal.Config
	log        *zap.Logger
	repos      *repository.Repositories
	evm        *rpc.EVMClient
	events     *events.NATSClient
	subscriber *Subscriber
	fetcher    *Fetcher
	canonical  *Canonical
	persist    *Persistence
	backfill   *Backfill
	head       atomic.Uint64
}

func NewManager(cfg internal.Config, log *zap.Logger, repos *repository.Repositories, evm *rpc.EVMClient, events *events.NATSClient) *Manager {
	return &Manager{
		cfg:        cfg,
		log:        log,
		repos:      repos,
		evm:        evm,
		events:     events,
		subscriber: NewSubscriber(evm, log),
		fetcher:    NewFetcher(evm, log),
		canonical:  NewCanonical(log),
		persist:    NewPersistence(repos, log),
		backfill:   NewBackfill(cfg, log),
	}
}

func (m *Manager) Start(ctx context.Context) {
	refs := m.subscriber.Run(ctx)
	backfillRefs := make(chan BlockRef, 8)
	go m.backfill.Run(ctx, backfillRefs)

	merged := mergeRefs(refs, backfillRefs)
	envelopes := m.fetcher.Run(ctx, merged)
	canonical := m.canonical.Run(ctx, envelopes, &m.head)
	go m.persist.Run(ctx, canonical)
}

func mergeRefs(a, b <-chan BlockRef) <-chan BlockRef {
	out := make(chan BlockRef, 16)
	var wg sync.WaitGroup
	merge := func(ch <-chan BlockRef) {
		defer wg.Done()
		for ref := range ch {
			out <- ref
		}
	}
	wg.Add(2)
	go merge(a)
	go merge(b)
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (m *Manager) Head() models.Block {
	return models.Block{
		Number: m.head.Load(),
	}
}
