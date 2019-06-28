package main

import (
	"log"
	"net/http"

	"github.com/levady/gohealth/cmd/serverd/httphandlers/home"
)

func main() {
	http.HandleFunc("/", home.Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
