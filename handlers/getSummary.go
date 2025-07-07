package handlers

import (
	"fmt"
	"net/http"
)

func GetSummary(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetSummary")
}
