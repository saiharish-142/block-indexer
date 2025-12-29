package pipeline

import (
	"context"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/events"
	"github.com/example/block-indexer/pkg/models"
	"github.com/example/block-indexer/pkg/rpc"
	"github.com/example/block-indexer/services/dag-indexer/internal"
	"github.com/example/block-indexer/services/dag-indexer/internal/repository"
)

type Manager struct {
	cfg        internal.Config
	log        *zap.Logger
	repos      *repository.Repositories
	dag        *rpc.DAGClient
	events     *events.NATSClient
	subscriber *Subscriber
	fetcher    *Fetcher
	persist    *Persistence
	tipMonitor *TipMonitor
	head       atomic.Uint64
}

func NewManager(cfg internal.Config, log *zap.Logger, repos *repository.Repositories, dag *rpc.DAGClient, events *events.NATSClient) *Manager {
	return &Manager{
		cfg:        cfg,
		log:        log,
		repos:      repos,
		dag:        dag,
		events:     events,
		subscriber: NewSubscriber(dag, log),
		fetcher:    NewFetcher(dag, log),
		persist:    NewPersistence(repos, log),
		tipMonitor: NewTipMonitor(dag, repos, log),
	}
}

func (m *Manager) Start(ctx context.Context) {
	refs := m.subscriber.Run(ctx)
	envelopes := m.fetcher.Run(ctx, refs)
	go m.persist.Run(ctx, envelopes, &m.head)
	go m.tipMonitor.Run(ctx)
}

func mergeRefs(a, b <-chan DagRef) <-chan DagRef {
	out := make(chan DagRef, 16)
	var wg sync.WaitGroup
	merge := func(ch <-chan DagRef) {
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

func (m *Manager) Head() models.DagBlock {
	return models.DagBlock{
		Order: m.head.Load(),
	}
}
