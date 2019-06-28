package httphandlers

import (
	"html/template"
	"net/http"

	"github.com/levady/gohealth/internal/sitehealthchecker"
)

// SiteHealthHandler represents SiteHealthHandler data
type SiteHealthHandler struct {
	Checker *sitehealthchecker.SiteHealthChecker
}

// Homepage renders the home page
func (handler *SiteHealthHandler) Homepage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/homepage.html")
	t.Execute(w, nil)
}
