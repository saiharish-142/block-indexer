package handlers

import (
	"net/http"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
)

func Search(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		results, err := deps.Search.Search(r.Context(), q)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": results})
	}
}
