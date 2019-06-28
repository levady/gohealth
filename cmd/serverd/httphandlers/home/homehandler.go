package home

import (
	"fmt"
	"net/http"
)

// Index renders the home page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, this is a home page")
}
