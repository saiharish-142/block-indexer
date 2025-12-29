package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
)

func GetTrace(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		trace, err := deps.Trace.Trace(r.Context(), hash)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, trace)
	}
}

func GetStateDiff(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		diff, err := deps.Trace.StateDiff(r.Context(), hash)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": diff})
	}
}
