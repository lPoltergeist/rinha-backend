package handlers

import (
	"encoding/json"
	"net/http"

	worker "github.com/lPoltergeist/rinha-backend.git/workers"
)

func GetSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := worker.Summary.BuildSummary()

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
	}
}
