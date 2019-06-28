package home

import (
	"html/template"
	"net/http"
)

// Index renders the home page
func Index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/templates/home/index.html")
	t.Execute(w, nil)
}
