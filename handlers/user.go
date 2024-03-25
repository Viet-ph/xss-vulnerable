package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/Viet-ph/xss-vulnerable/database"
	"github.com/Viet-ph/xss-vulnerable/models"
	"github.com/Viet-ph/xss-vulnerable/response"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	req := requestBody{}
	err := decoder.Decode(&req)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	for _, v := range db.Users {
		if v.Email == req.Email {
			response.RespondWithError(w, http.StatusForbidden, "Email already exist")
			return
		}
	}

	id := len(db.Users) + 1
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %s", err)
		return
	}

	newUser := models.User{Id: id, Email: req.Email, Password: string(hashedPass)}
	res := responseBody{Id: id, Email: req.Email}
	db.Users[id] = newUser

	err = database.Db.WriteDB(db)
	if err != nil {
		log.Printf("Error write DB to disk: %s", err)
		return
	}

	response.RespondWithJSON(w, 201, res)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	returnSlice := make([]models.User, 0, len(db.Users))
	for _, user := range db.Users {
		returnSlice = append(returnSlice, user)
	}

	sort.Slice(returnSlice, func(i, j int) bool { return returnSlice[i].Id < returnSlice[j].Id })

	response.RespondWithJSON(w, 200, returnSlice)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		log.Printf("Error convert string wildcard to int: %s", err)
		return
	}

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	user, exist := db.Users[userId]
	if exist {
		response.RespondWithJSON(w, 200, user)
	} else {
		response.RespondWithError(w, http.StatusNotFound, "Cannot find Chirpy with associated ID")
	}
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}

	jwtToken, _ := r.Context().Value("token").(*jwt.Token)
	userId, _ := ExtractIdFromToken(jwtToken)

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	id, _ := strconv.Atoi(userId)
	for i, v := range db.Users {
		if v.Id == id {
			decoder := json.NewDecoder(r.Body)
			req := requestBody{}
			err := decoder.Decode(&req)

			if err != nil {
				// an error will be thrown if the JSON is invalid or has the wrong types
				// any missing fields will simply have their values in the struct set to their zero value
				log.Printf("Error decoding parameters: %s", err)
				response.RespondWithError(w, http.StatusInternalServerError, "Error decoding request")
				return
			}

			newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Password hashing error: %s", err)
				response.RespondWithError(w, http.StatusInternalServerError, "New password hashing error")
				return
			}
			v.Email = req.Email
			v.Password = string(newHashedPassword)

			db.Users[i] = v

			err = database.Db.WriteDB(db)
			if err != nil {
				log.Printf("Error write DB to disk: %s", err)
				response.RespondWithError(w, http.StatusInternalServerError, "Error write DB to disk")
				return
			}

			response.RespondWithJSON(w, http.StatusOK, responseBody{Email: v.Email, Id: v.Id})
			return
		}
	}
	response.RespondWithError(w, http.StatusNotFound, "")
}
