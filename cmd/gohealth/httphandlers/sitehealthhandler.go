package httphandlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/levady/gohealth/internal/platform/sitestore"
	"github.com/levady/gohealth/internal/sitehealthchecker"
)

// Payload represents data to be displayed in the HTML template
type Payload struct {
	Data      interface{}
	ErrorData interface{}
}

// SiteHealthHandler represents SiteHealthHandler data
type SiteHealthHandler struct {
	SiteStore           *sitestore.Store
	HealtchCheckTimeout time.Duration
}

var runHealthChecks = healthChecksMethod
var homepageTplPath = "web/templates/homepage.html"

// Homepage renders the home page
func (handler *SiteHealthHandler) Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	p := Payload{Data: handler.SiteStore.List()}
	renderHomepage(w, p, http.StatusOK)
}

// Save saves a site to the store
func (handler *SiteHealthHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	url := r.FormValue("url")
	s := sitestore.Site{URL: strings.TrimSpace(url), Healthy: nil}

	if err := handler.SiteStore.Add(s); err != nil {
		errData := struct{ Msg string }{err.Error()}
		p := Payload{Data: handler.SiteStore.List(), ErrorData: errData}
		renderHomepage(w, p, http.StatusUnprocessableEntity)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// HealthChecks execute health checks on all stored sites
func (handler *SiteHealthHandler) HealthChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	runHealthChecks(handler.SiteStore, handler.HealtchCheckTimeout)

	json, err := json.Marshal(handler.SiteStore.List())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func renderHomepage(w http.ResponseWriter, p Payload, statusCode int) error {
	t, _ := template.ParseFiles(homepageTplPath)
	w.WriteHeader(statusCode)
	return t.Execute(w, p)
}

func healthChecksMethod(str *sitestore.Store, to time.Duration) {
	sitehealthchecker.ParallelHealthChecks(str, to)
}
