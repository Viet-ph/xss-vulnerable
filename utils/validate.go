package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CommentValidateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnCleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		fmt.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Comment is too long")
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	bodyDat := params.Body
	bodyDatSplit := strings.Split(bodyDat, " ")
	for _, profane := range profaneWords {
		for i, word := range bodyDatSplit {
			if strings.ToLower(word) == profane {
				bodyDatSplit[i] = "****"
			}
		}
	}

	bodyDat = strings.Join(bodyDatSplit, " ")
	RespondWithJSON(w, http.StatusOK, returnCleanedBody{CleanedBody: bodyDat})
}
