package main

import (
	"errors"
	"flag"
	"github.com/benjamin-vq/chirpy/internal/assert"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/benjamin-vq/chirpy/internal/database"
)

const (
	port = ":8080"

	fsDir = "."

	dbFilename = "database.json"

	fsPath            = "/app/*"
	readinessPath     = "GET /api/healthz"
	metricsPath       = "GET /admin/metrics"
	resetMetricsPath  = "GET /api/reset"
	postChirpPath     = "POST /api/chirps"
	getChirpsPath     = "GET /api/chirps"
	getChirpIdPath    = "GET /api/chirps/{chirpId}"
	postUsersPath     = "POST /api/users"
	loginPath         = "POST /api/login"
	putUsersPath      = "PUT /api/users"
	postRefreshPath   = "POST /api/refresh"
	postRevokePath    = "POST /api/revoke"
	deleteChirpIdPath = "DELETE /api/chirps/{chirpId}"
)

var debug = flag.Bool("debug", false, "Start on debug mode")

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func setupFlags() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	if debug != nil && *debug {
		log.Printf("[DEBUG] Deleting database file to start with a fresh one")
		err := os.Remove(dbFilename)
		assert.That(err == nil || errors.Is(err, os.ErrNotExist), "[DEBUG] Could not delete database file: %q", err)
	}

	err := godotenv.Load()
	assert.NoError(err, "Could not load environment variables")
}
func main() {
	setupFlags()

	db, err := database.NewDB(dbFilename)

	if err != nil {
		log.Fatalf("Error creating database: %q", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	assert.That(jwtSecret != "", "Jwt Secret should not be empty")

	apiConfig := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()

	fileserverHandler := http.StripPrefix("/app", http.FileServer(http.Dir(fsDir)))
	mux.Handle(fsPath, apiConfig.metricsIncrementer(fileserverHandler))

	mux.HandleFunc(readinessPath, readinessHandler)
	mux.HandleFunc(metricsPath, apiConfig.metricsHandler)
	mux.HandleFunc(resetMetricsPath, apiConfig.metricsReseter)
	mux.HandleFunc(postChirpPath, apiConfig.postChirpHandler)
	mux.HandleFunc(getChirpsPath, apiConfig.getChirpHandler)
	mux.HandleFunc(getChirpIdPath, apiConfig.chirpIdGetHandler)
	mux.HandleFunc(postUsersPath, apiConfig.postUsersHandler)
	mux.HandleFunc(loginPath, apiConfig.loginPostHandler)
	mux.HandleFunc(putUsersPath, apiConfig.putUsersHandler)
	mux.HandleFunc(postRefreshPath, apiConfig.postRefreshHandler)
	mux.HandleFunc(postRevokePath, apiConfig.postRevokeHandler)
	mux.HandleFunc(deleteChirpIdPath, apiConfig.deleteChirpIdHandler)

	log.Printf("Registered file handler for dir %q on path %q", fsDir, fsPath)
	log.Printf("Registered readiness endpoint on path %q", readinessPath)
	log.Printf("Registered metrics endpoint on path %q", metricsPath)
	log.Printf("Registered reset metrics endpoint on path %q", resetMetricsPath)
	log.Printf("Registered POST chirps endpoint on path %q", postChirpPath)
	log.Printf("Registered GET chirps endpoint on path %q", getChirpsPath)
	log.Printf("Registered GET chirp by id endpoint on path %q", getChirpIdPath)
	log.Printf("Registered POST users endpoint on path %q", postUsersPath)
	log.Printf("Registered PUT users endpoint on path %q", putUsersPath)
	log.Printf("Registered POST login endpoint on path %q", loginPath)
	log.Printf("Registered POST refresh endpoint on path %q", postRefreshPath)
	log.Printf("Registered POST revoke endpoint on path %q", postRevokePath)
	log.Printf("Registered DELETE chirp by id endpoint on path %q", deleteChirpIdPath)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Starting server on port %q", port)
	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("error starting server: %q", err)
	}
}
