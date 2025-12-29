package handlers

import (
	"net/http"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
)

func GetLogs(_ *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{}, "query": r.URL.RawQuery})
	}
}
