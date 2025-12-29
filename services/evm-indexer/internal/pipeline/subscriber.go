package pipeline

import (
	"context"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/rpc"
)

type BlockRef struct {
	Hash   string
	Number uint64
}

type Subscriber struct {
	client *rpc.EVMClient
	log    *zap.Logger
}

func NewSubscriber(client *rpc.EVMClient, log *zap.Logger) *Subscriber {
	return &Subscriber{client: client, log: log}
}

func (s *Subscriber) Run(ctx context.Context) <-chan BlockRef {
	out := make(chan BlockRef, 8)
	heads, err := s.client.SubscribeNewHeads(ctx)
	if err != nil {
		s.log.Warn("subscribe new heads failed, falling back to polling", zap.Error(err))
		close(out)
		return out
	}
	go func() {
		defer close(out)
		for head := range heads {
			out <- BlockRef{Hash: head.Hash, Number: head.Number}
		}
	}()
	return out
}
