package main

import (
	"net/http"

	"github.com/lPoltergeist/rinha-backend.git/handlers"
	"github.com/lPoltergeist/rinha-backend.git/helper"
)

var HealthCheck bool = true

func main() {
	helper.StartHealthCheck()
	handlers.InitWorkers(30)

	http.HandleFunc("/", handlers.HelloWorld)
	http.HandleFunc("/payments", handlers.Payments)
	http.HandleFunc("/payments-summary", handlers.GetSummary)

	http.ListenAndServe(":9999", nil)
}
