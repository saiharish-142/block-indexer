package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/block-indexer/services/stats-service/internal/repository"
)

func Overview(repo *repository.StatsRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := repo.Overview(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(stats)
	}
}

func Historical(repo *repository.StatsRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := r.URL.Query().Get("metric")
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")
		var from, to time.Time
		if fromStr != "" {
			from, _ = time.Parse(time.RFC3339, fromStr)
		}
		if toStr != "" {
			to, _ = time.Parse(time.RFC3339, toStr)
		}
		points, err := repo.Historical(r.Context(), metric, from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"items": points})
	}
}
