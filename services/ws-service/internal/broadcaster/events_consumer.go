package broadcaster

import (
	"context"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/events"
	"github.com/example/block-indexer/services/ws-service/internal/hub"
)

type EventsConsumer struct {
	log    *zap.Logger
	hub    *hub.Hub
	events *events.NATSClient
}

func New(log *zap.Logger, h *hub.Hub, events *events.NATSClient) *EventsConsumer {
	return &EventsConsumer{log: log, hub: h, events: events}
}

func (c *EventsConsumer) Start(ctx context.Context) {
	if c.events == nil {
		return
	}
	handler := func(msg []byte) error {
		c.hub.Broadcast(hub.Message{Type: "event", Data: string(msg)})
		return nil
	}
	if _, err := c.events.Subscribe(ctx, events.TopicEvmBlockCanonical, handler); err != nil {
		c.log.Warn("subscribe failed", zap.Error(err))
	}
}
