package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) metricsReseter(w http.ResponseWriter, req *http.Request) {
	log.Println("Resetting metrics back to 0")

	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusNoContent)
}
