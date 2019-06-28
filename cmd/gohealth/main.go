package main

import (
	"log"
	"net/http"

	"github.com/levady/gohealth/cmd/gohealth/httphandlers"
)

func main() {
	http.HandleFunc("/", httphandlers.Homepage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
