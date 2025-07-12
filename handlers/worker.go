package handlers

import "github.com/lPoltergeist/rinha-backend.git/models"

var JobChan chan models.Payment

func InitWorkers(n int) {
	JobChan = make(chan models.Payment, 1000)

	for i := 0; i < n; i++ {
		go worker(i, JobChan)
	}
}
