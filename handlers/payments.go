package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lPoltergeist/rinha-backend.git/models"
	"github.com/lPoltergeist/rinha-backend.git/queue"
)

func Payments(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var payment models.Payment

	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if payment.RequestAt == "" {
		payment.RequestAt = time.Now().UTC().Format(time.RFC3339)
	}

	if err := queue.Enqueued(payment); err != nil {
		fmt.Fprint(w, "Error to queue your payment!", http.StatusInternalServerError)
	}

	fmt.Printf("time to complete the handle: %v\n", time.Since(start))
	w.WriteHeader(http.StatusAccepted)

}
