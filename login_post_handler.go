package main

import (
	"encoding/json"
	"github.com/benjamin-vq/chirpy/internal/auth"
	"log"
	"net/http"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User
}

func (cfg *apiConfig) loginPostHandler(w http.ResponseWriter, r *http.Request) {

	loginParams := LoginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginParams)

	if err != nil {
		log.Printf("Error decoding login parameters: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not decode login parameters")
		return
	}

	user, err := cfg.DB.UserByEmail(loginParams.Email)

	if err != nil {
		log.Printf("Error finding user by email: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	err = auth.ComparePasswordHash(user.Password, loginParams.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			user.Email,
			user.Id,
		}})
}
