package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/benjamin-vq/chirpy/internal/assert"
)

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	errorResponse := errorResponse{}

	errorResponse.Error = msg
	bytes, err := json.Marshal(&errorResponse)

	if err != nil {
		log.Printf("Error mashalling error response: %q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	assert.That(code < 400, "Code should be in the 100-399 range")

	bytes, err := json.Marshal(&payload)

	if err != nil {
		log.Printf("Error mashalling chirp: %q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bytes)
}
