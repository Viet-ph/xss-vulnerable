package handlers

import (
	"net/http"
	"text/template"
)

func SocialHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/social.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	tmpl.Execute(w, nil)
}
