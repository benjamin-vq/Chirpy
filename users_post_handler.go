package main

import (
	"encoding/json"
	"errors"
	"github.com/benjamin-vq/chirpy/internal/auth"
	"github.com/benjamin-vq/chirpy/internal/database"
	"log"
	"net/http"
)

type User struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type userParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) postUsersHandler(w http.ResponseWriter, r *http.Request) {

	params := userParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding user: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not decode user")
		return
	}

	validParams, err := validateParams(params)

	if err != nil {
		log.Printf("Validation of user parameters failed: %q", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := cfg.DB.CreateUser(validParams.Email, validParams.Password)
	if err != nil {
		if errors.Is(err, database.ErrEmailExists) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("Could not save user to database: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		user.Email,
		user.Id,
		user.IsChirpyRed,
	})
}

func validateParams(p userParams) (userParams, error) {

	if p.Email == "" {
		log.Print("Received an empty email during validation, returning error")
		return userParams{}, errors.New("email can not be empty")
	}

	if p.Password == "" {
		log.Print("Received an empty password during validation, returning error")
		return userParams{}, errors.New("password can not be empty")
	}

	hashed, err := auth.HashPassword(p.Password)
	if err != nil {
		log.Printf("Could not generate hash from password: %q", err)
		return userParams{}, err
	}

	return userParams{Email: p.Email, Password: hashed}, nil
}
