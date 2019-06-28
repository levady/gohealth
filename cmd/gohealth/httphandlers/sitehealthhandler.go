package httphandlers

import (
	"html/template"
	"net/http"
)

// Homepage renders the home page
func Homepage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/homepage.html")
	t.Execute(w, nil)
}
