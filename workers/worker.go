package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lPoltergeist/rinha-backend.git/models"
	"github.com/lPoltergeist/rinha-backend.git/queue"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func InitWorkers(n int) {
	paymentChan := make(chan models.Payment, 2000)

	go func() {
		ch := queue.Consumer()

		for payment := range ch {
			paymentChan <- payment
		}
	}()

	for i := 0; i < n; i++ {
		go func(id int) {
			for payment := range paymentChan {
				if err := Start(payment); err != nil {
					fmt.Printf("Worker %d failed to process payment: %v\n", id, err)
				}

			}
		}(i)
	}
}

func Start(payment models.Payment) error {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	jsonBody, err := json.Marshal(payment)

	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}

	baseSleepTime := 10

	urls := []string{
		"http://payment-processor-default:8080/payments",
		"http://payment-processor-fallback:8080/payments",
	}

	for attempts := 0; attempts <= 10; attempts++ {
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

			elapse := time.Since(start)

			fmt.Printf("delay til 200: %v\n", elapse)
			Summary.Add(payment.Amount, processor)

			return nil
		}

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}

	return nil
}
