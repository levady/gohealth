package httphandlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/levady/gohealth/internal/platform/sitestore"
)

// Payload represents data to be displayed in the HTML template
type Payload struct {
	Data      interface{}
	ErrorData interface{}
}

// SiteHealthHandler represents SiteHealthHandler data
type SiteHealthHandler struct {
	SiteStore *sitestore.Store
}

var homepageTplPath = "web/templates/homepage.html"

// Homepage renders the home page
func (handler *SiteHealthHandler) Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	p := Payload{Data: handler.SiteStore.List()}
	renderHomepage(w, p, http.StatusOK)
}

// Save saves a site to the store
func (handler *SiteHealthHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
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

// Delete deletes a site from the store
func (handler *SiteHealthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.NotFound(w, r)
		return
	}

	siteIDStr := r.URL.Path[len("/ajax/sites/delete/"):]
	siteID, err := strconv.ParseInt(siteIDStr, 10, 64)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	if err := handler.SiteStore.Delete(siteID); err != nil {
		http.NotFound(w, r)
		return
	}

	// Must return JSON response, if not jQuery won't fire the `success` callback and
	// fire the error callback instead....
	json, err := json.Marshal(struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// HealthChecks execute health checks on all stored sites
func (handler *SiteHealthHandler) HealthChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

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
