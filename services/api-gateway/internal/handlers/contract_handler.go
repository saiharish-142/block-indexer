package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/block-indexer/services/api-gateway/internal/clients"
)

func GetContract(deps *clients.ClientSet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := chi.URLParam(r, "address")
		contract, err := deps.Contract.Get(r.Context(), addr)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, contract)
	}
}
