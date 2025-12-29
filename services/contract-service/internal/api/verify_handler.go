package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/block-indexer/pkg/models"
	"github.com/example/block-indexer/services/contract-service/internal/repository"
	"github.com/example/block-indexer/services/contract-service/internal/verifier"
)

type verifyRequest struct {
	Source   string         `json:"source"`
	Metadata map[string]any `json:"metadata"`
}

func VerifyContract(repo *repository.ContractsRepo, compiler verifier.Compiler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := chi.URLParam(r, "addr")
		var req verifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		bytecode, err := compiler.Compile(verifier.CompileRequest{Source: req.Source, Settings: req.Metadata})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		contract := models.Contract{
			Address:  addr,
			Verified: true,
			ABI:      req.Metadata,
			Bytecode: bytecode,
		}
		if err := repo.Upsert(r.Context(), contract); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(contract)
	}
}
