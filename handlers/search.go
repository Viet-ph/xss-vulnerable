package handlers

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Viet-ph/xss-vulnerable/database"
	"github.com/golang-jwt/jwt/v5"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var userEmail string
	jwtToken, _ := r.Context().Value("token").(*jwt.Token)
	userId, _ := ExtractIdFromToken(jwtToken)

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	id, _ := strconv.Atoi(userId)
	for _, v := range db.Users {
		if v.Id == id {
			userEmail = v.Email
		}
	}

	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Directly taking user input from query parameters
	userQuery := r.URL.Query().Get("query")

	// Render template with user input - this is vulnerable to XSS!
	err = tmpl.Execute(w, map[string]interface{}{
		"Query": userQuery,
		"Email": userEmail,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}
