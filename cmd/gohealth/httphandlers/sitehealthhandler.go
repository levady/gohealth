package httphandlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/levady/gohealth/internal/sitehealthchecker"
)

// Payload represents data to be displayed in the HTML template
type Payload struct {
	Data      interface{}
	ErrorData interface{}
}

// SiteHealthHandler represents SiteHealthHandler data
type SiteHealthHandler struct {
	Checker *sitehealthchecker.SiteHealthChecker
}

// Homepage renders the home page
func (handler *SiteHealthHandler) Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	p := Payload{Data: handler.Checker.Sites}
	renderHomepage(w, p)
}

// Save saves a site to the store
func (handler *SiteHealthHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	url := r.FormValue("url")
	s := sitehealthchecker.Site{URL: url, Healthy: nil}

	if err := handler.Checker.AddSite(s); err != nil {
		errData := struct{ Msg string }{err.Error()}
		p := Payload{Data: handler.Checker.Sites, ErrorData: errData}
		renderHomepage(w, p)
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

	handler.Checker.RunHealthChecks()

	json, err := json.Marshal(handler.Checker.Sites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func renderHomepage(w http.ResponseWriter, p Payload) {
	t, _ := template.ParseFiles("web/templates/homepage.html")
	t.Execute(w, p)
}
