package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	//"github.com/lPoltergeist/rinha-backend.git/helper"
	"github.com/lPoltergeist/rinha-backend.git/models"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
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
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	jsonBody, err := json.Marshal(payment)

	if err != nil {
		fmt.Println(err)
	}

	baseSleepTime := 100

	urls := []string{
		"http://payment-processor-default:8080/payments",
		"http://payment-processor-fallback:8080/payments",
	}

	for attempts := 0; attempts <= 100; attempts++ {
		url := urls[attempts%2]
		sleepTime := baseSleepTime * (1 << attempts)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		defer res.Body.Close()
		io.Copy(io.Discard, res.Body)

		if res.StatusCode == 200 {
			processor := false
			if strings.Contains(url, "fallback") {
				processor = true
			}

			Summary.Add(payment.Amount, processor)

			return res.Status, nil
		}

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
