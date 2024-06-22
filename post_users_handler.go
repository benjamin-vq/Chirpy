package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) postUsersHandler(w http.ResponseWriter, r *http.Request) {

	params := struct {
		Email string `json:"email"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding user: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not decode user")
		return
	}

	if params.Email == "" {
		log.Printf("Received email for user was empty, responding with error")
		respondWithError(w, http.StatusBadRequest, "Email can not be empty")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		log.Printf("Could not save user to database: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
