package handlers

import (
	"encoding/json"
	"net/http"
)

func GetSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Summary.BuildSummary()

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erro ao gerar JSON", http.StatusInternalServerError)
	}
}
