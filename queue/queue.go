package queue

import (
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

func Consumer() <-chan models.Payment {
	ch := make(chan models.Payment)

	go func() {
		for {
			if data.Client == nil {
				fmt.Printf("client is nil %v\n", data.Client)
			}

			res, err := data.Client.BRPop(data.Context, 0, "payment").Result()
			if err != nil {
				fmt.Printf("Error on consumer, fix it! %v", err)
				continue
			}

			var payment models.Payment

			if err := json.Unmarshal([]byte(res[1]), &payment); err == nil {
				ch <- payment
			} else {
				fmt.Printf("Error on consumer, fix it! %v", err)
			}
		}
	}()

	return ch
}
