package handlers

import (
	"fmt"
	"net/http"
)

func (h *httpHandler)HandleHome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello!\n")
}
