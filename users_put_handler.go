package main

import (
	"encoding/json"
	"github.com/benjamin-vq/chirpy/internal/auth"
	"github.com/benjamin-vq/chirpy/internal/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type updateParams struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (cfg *apiConfig) putUsersHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	token, found := strings.CutPrefix(authHeader, "Bearer ")

	if !found || token == "" {
		log.Printf("Did not find a valid authorization header")
		respondWithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	subject, err := auth.ValidateToken(token, cfg.jwtSecret)

	if err != nil {
		log.Printf("Token validation returned an error: %q", err)
		respondWithError(w, http.StatusUnauthorized, "Could not validate token")
		return
	}

	params := updateParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding update params: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not decode user information")
		return
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		log.Printf("Decoded subject %s is not a valid user id", subject)
		respondWithError(w, http.StatusInternalServerError, "Invalid user id")
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Could not hash new password: %q", err)
		respondWithError(w, http.StatusInternalServerError, "An internal error occurred.")
	}

	err = cfg.DB.UpdateUser(&database.User{params.Email, newHashedPassword, id})
	if err != nil {
		log.Printf("Could not update user: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	type response struct {
		User
	}
	log.Printf("Succesfully update user")
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Email: params.Email,
			ID:    id,
		},
	})
}
