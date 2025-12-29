package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
)

func GetTransaction(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		tx, err := deps.EVM.GetTransaction(r.Context(), hash)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, tx)
	}
}

func GetAddressTxs(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		items, err := deps.EVM.GetAddressTxs(r.Context(), address)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}
