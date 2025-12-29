package pipeline

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/example/block-indexer/pkg/models"
	"github.com/example/block-indexer/pkg/rpc"
	"github.com/example/block-indexer/services/dag-indexer/internal/repository"
)

type TipMonitor struct {
	client *rpc.DAGClient
	repos  *repository.Repositories
	log    *zap.Logger
}

func NewTipMonitor(client *rpc.DAGClient, repos *repository.Repositories, log *zap.Logger) *TipMonitor {
	return &TipMonitor{client: client, repos: repos, log: log}
}

func (m *TipMonitor) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tips, err := m.client.GetTips(ctx)
			if err != nil {
				m.log.Warn("failed to fetch tips", zap.Error(err))
				continue
			}

			// Convert to models
			var modelTips []models.BlockTip
			for _, t := range tips.Valid {
				modelTips = append(modelTips, models.BlockTip{
					Hash:   t.Hash,
					Height: t.Height,
					Status: "valid",
				})
			}
			for _, t := range tips.Invalid {
				modelTips = append(modelTips, models.BlockTip{
					Hash:   t.Hash,
					Height: t.Height,
					Status: "invalid",
				})
			}

			if err := m.repos.Tips.UpdateTips(ctx, modelTips); err != nil {
				m.log.Warn("failed to update tips", zap.Error(err))
			}
		}
	}
}
