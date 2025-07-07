package main

import (
	"net/http"

	"github.com/lPoltergeist/rinha-backend.git/handlers"
)

func main() {

	http.HandleFunc("/", handlers.HelloWorld)
	http.HandleFunc("/payments", handlers.Payments)
	http.HandleFunc("/payments-summary", handlers.GetSummary)

	http.ListenAndServe(":9999", nil)
}
