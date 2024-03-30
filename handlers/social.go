package handlers

import (
	"net/http"
	"text/template"

	"github.com/golang-jwt/jwt/v5"
)

func SocialHandler(w http.ResponseWriter, r *http.Request) {
	jwtToken, _ := r.Context().Value("token").(*jwt.Token)
	userEmail, err := ExtractEmailFromToken(jwtToken)
	if err != nil{
		http.Error(w, "Internal Server Error", 500)
		return
	}

	tmpl, err := template.ParseFiles("templates/social.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Render template with user input - this is vulnerable to XSS!
	err = tmpl.Execute(w, map[string]interface{}{
		"Email": userEmail,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}
