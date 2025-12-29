package api

import (
	"encoding/json"
	"net/http"

	"github.com/example/block-indexer/services/search-service/internal/repository"
)

func Search(repo *repository.SearchRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		results, err := repo.Search(r.Context(), q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"items": results})
	}
}
