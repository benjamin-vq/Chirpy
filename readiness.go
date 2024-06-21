package main

import "net/http"

func readinessHandler(w http.ResponseWriter, req *http.Request) {

	bytes := []byte("OK")

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(bytes)

}
