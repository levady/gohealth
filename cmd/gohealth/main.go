package main

import (
	"log"
	"net/http"
	"time"

	"github.com/levady/gohealth/cmd/gohealth/httphandlers"
	"github.com/levady/gohealth/internal/sitehealthchecker"
)

func main() {
	shc := sitehealthchecker.New(15 * time.Second)
	shh := httphandlers.SiteHealthHandler{
		Checker: &shc,
	}

	http.HandleFunc("/", shh.Homepage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
