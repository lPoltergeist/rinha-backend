package handlers

import (
	"fmt"
	"net/http"
)

func Payments(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pagamentos")
}
