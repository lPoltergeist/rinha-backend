package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lPoltergeist/rinha-backend.git/models"
)

func sendToPaymentProcessor(payment models.Payment) (string, error) {
	jsonBody, err := json.Marshal(payment)

	if err != nil {
		fmt.Println(err)
	}

	attempts := 0

	urls := []string{
		"http://192.168.1.109:8001/payments",
		"http://192.168.1.109:8002/payments",
	}

	baseSleepTime := 100

	for {
		url := urls[attempts%2]
		sleepTime := baseSleepTime * (2 ^ attempts)

		res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			attempts++
			continue
		}

		fmt.Printf("Response: %v\n", res.Status)

		if res.StatusCode == 200 {
			fmt.Printf("Success Payment: %v\n", res.Status)
			return res.Status, nil
		}

		fmt.Printf("Attempts: %v\n", attempts)

		attempts++

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}
func Payments(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment

	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if payment.RequestAt == "" {
		payment.RequestAt = time.Now().UTC().Format(time.RFC3339)
	}

	response, _ := sendToPaymentProcessor(payment)

	fmt.Fprintf(w, "Pagamentos: %v", response)
}
