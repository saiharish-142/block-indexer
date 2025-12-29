package handlers

import (
	"net/http"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
)

func GetStats(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := deps.Stats.Overview(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, stats)
	}
}
