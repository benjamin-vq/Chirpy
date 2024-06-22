package main

import (
	"encoding/json"
	"errors"
	"github.com/benjamin-vq/chirpy/internal/auth"
	"github.com/benjamin-vq/chirpy/internal/database"
	"log"
	"net/http"
	"time"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Expires  int    `json:"expires_in_seconds"`
}

type LoginResponse struct {
	User
	Token string `json:"token"`
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
		if errors.Is(err, database.UserNotExists) {
			log.Printf("User with email %q does not exist %q", loginParams.Email, err)
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		log.Printf("Error finding user by email: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	err = auth.ComparePasswordHash(user.Password, loginParams.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jwt, err := auth.CreateJwt(time.Duration(loginParams.Expires), user.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			user.Email,
			user.Id,
		},
		Token: jwt,
	})
}
