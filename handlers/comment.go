package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"text/template"

	"github.com/Viet-ph/xss-vulnerable/database"
	"github.com/Viet-ph/xss-vulnerable/models"
	"github.com/Viet-ph/xss-vulnerable/utils"
	"github.com/golang-jwt/jwt/v5"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(req.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	id := len(db.Comments) + 1
	newComment := models.Comment{Id: id, Body: req.Body}
	db.Comments[id] = newComment
	err = database.Db.WriteDB(db)
	if err != nil {
		log.Printf("Error write DB to disk: %s", err)
		return
	}
	utils.RespondWithJSON(w, 201, newComment)
}

func GetAllCommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	returnSlice := make([]models.Comment, 0, len(db.Comments))
	//returnSlice := []models.Chirp{}
	for _, chirpy := range db.Comments {
		returnSlice = append(returnSlice, chirpy)
	}

	sort.Slice(returnSlice, func(i, j int) bool { return returnSlice[i].Id < returnSlice[j].Id })

	utils.RespondWithJSON(w, 200, returnSlice)
}

func GetCommentById(w http.ResponseWriter, r *http.Request) {
	chirpId, err := strconv.Atoi(r.PathValue("commentID"))
	if err != nil {
		log.Printf("Error convert string wildcard to int: %s", err)
		return
	}

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	chirpy, exist := db.Comments[chirpId]
	if exist {
		utils.RespondWithJSON(w, 200, chirpy)
	} else {
		utils.RespondWithError(w, http.StatusNotFound, "Cannot find Chirpy with associated ID")
	}
}

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/comment.html")
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

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

	// Handle new comment submission
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
		comment := r.FormValue("comment")
		if comment != "" {
			id := len(db.Comments) + 1
			newComment := models.Comment{Id: id, Body: comment}
			db.Comments[id] = newComment
			err = database.Db.WriteDB(db)
			if err != nil {
				log.Printf("Error write DB to disk: %s", err)
				return
			}
		}
	}
	// Render template with comments
	var comments []string
	for _, v := range db.Comments {
		comments = append(comments, v.Body)
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"Email": userEmail, // Replace with actual username retrieval logic
		"Comments": comments,
	})

	if err != nil {
		fmt.Print(err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
	}
}
