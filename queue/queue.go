package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lPoltergeist/rinha-backend.git/data"
	"github.com/lPoltergeist/rinha-backend.git/models"
)

func Enqueued(payment models.Payment) error {
	json, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	return data.Client.LPush(data.Context, "payment", json).Err()
}

func Consumer(ctx context.Context) <-chan models.Payment {
	ch := make(chan models.Payment)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Consumer context done, exiting...")
				return
			default:
				// Continue the normal processing
			}

			var payment models.Payment

			if data.Client == nil {
				fmt.Printf("Client is nil")
				continue
			}

			res, err := data.Client.BRPop(data.Context, 0, "payment").Result()

			if err != nil {
				fmt.Printf("Error on getting result: %v\n", err)
				continue
			}

			if err := json.Unmarshal([]byte(res[1]), &payment); err != nil {
				fmt.Printf("Error on unmarshalling json %v\n", err)
				continue
			}

			select {
			case <-ctx.Done():
				fmt.Println("Consumer context done while sending payment, exiting...")
				return
			case ch <- payment:
			}
		}
	}()

	return ch
}
