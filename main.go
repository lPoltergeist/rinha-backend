package main

import (
	"net/http"

	"github.com/lPoltergeist/rinha-backend.git/data"
	"github.com/lPoltergeist/rinha-backend.git/handlers"
)

var HealthCheck bool = true

func main() {
	data.InitRedis()

	handlers.InitWorkers(190)

	http.HandleFunc("/", handlers.HelloWorld)
	http.HandleFunc("/payments", handlers.Payments)
	http.HandleFunc("/payments-summary", handlers.GetSummary)

	http.ListenAndServe(":9999", nil)
}
