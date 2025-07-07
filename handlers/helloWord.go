package handlers

import (
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, R *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
