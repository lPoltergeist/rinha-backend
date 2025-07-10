package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lPoltergeist/rinha-backend.git/models"
)

var client = &http.Client{
	Timeout: time.Second * 5,
}

func worker(id int, jobs <-chan models.Payment) {
	for payment := range jobs {
		_, err := sendToPaymentProcessor(payment)
		if nil != err {
			fmt.Printf("Worker %d failed: %v\n", id, err)
		}
	}
}

func sendToPaymentProcessor(payment models.Payment) (string, error) {
	jsonBody, err := json.Marshal(payment)

	if err != nil {
		fmt.Println(err)
	}

	urls := []string{
		"http://localhost:8001/payments",
		"http://localhost:8002/payments",
	}

	baseSleepTime := 100

	for attempts := 0; attempts <= 10; attempts++ {
		url := urls[attempts%2]
		sleepTime := baseSleepTime * (1 << attempts)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			attempts++
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			attempts++
			continue
		}

		defer res.Body.Close()
		io.Copy(io.Discard, res.Body)

		if res.StatusCode == 200 {
			fmt.Printf("Success Payment: %v\n", res.Status)
			return res.Status, nil
		}

		fmt.Printf("Attempts: %v\n", attempts)

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}

	return "", nil
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

	JobChan <- payment

	w.WriteHeader(http.StatusAccepted)
}
