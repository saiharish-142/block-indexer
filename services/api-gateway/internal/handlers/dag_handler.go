package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
)

func GetDagBlock(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		block, err := deps.DAG.GetBlock(r.Context(), hash)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, block)
	}
}

func GetDagRange(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, _ := strconv.ParseInt(r.URL.Query().Get("fromHeight"), 10, 64)
		to, _ := strconv.ParseInt(r.URL.Query().Get("toHeight"), 10, 64)
		items, err := deps.DAG.GetRange(r.Context(), from, to)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}
