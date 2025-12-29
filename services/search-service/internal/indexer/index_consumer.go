package indexer

import (
	"context"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/events"
	"github.com/example/block-indexer/services/search-service/internal/repository"
)

type IndexConsumer struct {
	log    *zap.Logger
	repo   *repository.SearchRepo
	events *events.NATSClient
}

func NewIndexConsumer(log *zap.Logger, repo *repository.SearchRepo, events *events.NATSClient) *IndexConsumer {
	return &IndexConsumer{log: log, repo: repo, events: events}
}

func (c *IndexConsumer) Start(ctx context.Context) {
	handler := func(msg []byte) error {
		c.log.Debug("index event", zap.ByteString("payload", msg))
		return nil
	}
	if c.events != nil {
		if _, err := c.events.Subscribe(ctx, events.TopicEvmTxInserted, handler); err == nil {
			c.log.Info("listening for tx inserted events")
		}
	}
}
