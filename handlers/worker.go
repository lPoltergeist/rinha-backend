package handlers

import "github.com/lPoltergeist/rinha-backend.git/models"

var JobChan chan models.Payment

func InitWorkers(n int) {
	JobChan = make(chan models.Payment, 100)

	for i := range n {
		go worker(i, JobChan)
	}
}
