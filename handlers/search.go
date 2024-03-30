package handlers

import (
	"net/http"
	"text/template"

	"github.com/golang-jwt/jwt/v5"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var userEmail string
	jwtToken, _ := r.Context().Value("token").(*jwt.Token)
	userEmail, err := ExtractEmailFromToken(jwtToken)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	
	// Directly taking user input from query parameters
	userQuery := r.URL.Query().Get("query")
	//userQuery = utils.InputValidater(userQuery)

	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Render template with user input - this is vulnerable to XSS!
	err = tmpl.Execute(w, map[string]interface{}{
		"Query": userQuery,
		"Email": userEmail,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}
