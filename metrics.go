package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {

	bytes := []byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, cfg.fileserverHits))
	w.Header().Add("Content-Type", "text/xml")
	w.WriteHeader(200)
	w.Write(bytes)
}

func (cfg *apiConfig) metricsIncrementer(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Incrementing hit metrics, current hits: %d", cfg.fileserverHits)

		cfg.fileserverHits += 1
		next.ServeHTTP(w, req)
	})

}
