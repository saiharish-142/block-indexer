package internal

import (
	"github.com/go-chi/chi/v5"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
	"github.com/example/block-indexer/services/api-gateway/internal/handlers"
)

func RegisterRoutes(router *chi.Mux, deps *clients.ClientSet) {
	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/blocks", handlers.ListBlocks(deps))
		r.Get("/blocks/{id}", handlers.GetBlock(deps))
		r.Get("/tx/{hash}", handlers.GetTransaction(deps))
		r.Get("/address/{address}/txs", handlers.GetAddressTxs(deps))
		r.Get("/logs", handlers.GetLogs(deps))
		r.Get("/dag/block/{hash}", handlers.GetDagBlock(deps))
		r.Get("/dag/graph", handlers.GetDagRange(deps))
		r.Get("/stats/overview", handlers.GetStats(deps))
		r.Get("/contracts/{address}", handlers.GetContract(deps))
		r.Get("/search", handlers.Search(deps))
		r.Get("/tx/{hash}/trace", handlers.GetTrace(deps))
		r.Get("/tx/{hash}/state-diff", handlers.GetStateDiff(deps))
	})
}
