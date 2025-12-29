package pipeline

import (
	"context"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/rpc"
)

type DagRef struct {
	Hash  string
	Order int64
}

type Subscriber struct {
	client *rpc.DAGClient
	log    *zap.Logger
}

func NewSubscriber(client *rpc.DAGClient, log *zap.Logger) *Subscriber {
	return &Subscriber{client: client, log: log}
}

func (s *Subscriber) Run(ctx context.Context) <-chan DagRef {
	out := make(chan DagRef, 8)
	stream, err := s.client.SubscribeDagBlocks(ctx)
	if err != nil {
		s.log.Warn("failed subscribing dag blocks", zap.Error(err))
		close(out)
		return out
	}
	go func() {
		defer close(out)
		for evt := range stream {
			out <- DagRef{Hash: evt.Hash, Order: evt.Order}
		}
	}()
	return out
}
