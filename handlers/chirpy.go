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
)

func CreateChirpyHandler(w http.ResponseWriter, r *http.Request) {
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
		response.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	id := len(db.Chirps) + 1
	newChirpy := models.Chirp{Id: id, Body: req.Body}
	db.Chirps[id] = newChirpy
	err = database.Db.WriteDB(db)
	if err != nil {
		log.Printf("Error write DB to disk: %s", err)
		return
	}
	response.RespondWithJSON(w, 201, newChirpy)
}

func GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	returnSlice := make([]models.Chirp, 0, len(db.Chirps))
	//returnSlice := []models.Chirp{}
	for _, chirpy := range db.Chirps {
		returnSlice = append(returnSlice, chirpy)
	}

	sort.Slice(returnSlice, func(i, j int) bool { return returnSlice[i].Id < returnSlice[j].Id })

	response.RespondWithJSON(w, 200, returnSlice)
}

func GetChirpyById(w http.ResponseWriter, r *http.Request) {
	chirpId, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error convert string wildcard to int: %s", err)
		return
	}

	db, err := database.Db.LoadDB()
	if err != nil {
		log.Printf("Error load DB to memory: %s", err)
		return
	}

	chirpy, exist := db.Chirps[chirpId]
	if exist {
		response.RespondWithJSON(w, 200, chirpy)
	} else {
		response.RespondWithError(w, http.StatusNotFound, "Cannot find Chirpy with associated ID")
	}

}
