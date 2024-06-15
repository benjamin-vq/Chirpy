package main

import (
	"log"
	"net/http"
)

const (
	port = ":8080"

	fsDir  = "."
	fsPath = "/app/*"

	readinessPath = "/healthz"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle(fsPath, http.StripPrefix("/app", http.FileServer(http.Dir(fsDir))))
	log.Printf("Registered file handler for dir %q on path %q", fsDir, fsPath)

	mux.HandleFunc(readinessPath, readinessHandler)
	log.Printf("Registered readiness endpoint on path: %q", readinessPath)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Starting server on port %q", port)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("error starting server: %q\n", err)
	}
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {

	bytes := []byte("OK")

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(bytes)

}
