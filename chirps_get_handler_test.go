package main

import (
	"errors"
	"fmt"
	"github.com/benjamin-vq/chirpy/internal/database"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestChirpsGetHandler(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := database.NewDB(testDbName)
	if err != nil {
		log.Fatalf("Could not create new database for test: %q", err)
	}

	cfg := apiConfig{
		DB: db,
	}

	empty := struct {
		wantCode int
		wantBody string
	}{
		wantCode: 204,
		wantBody: `[]`,
	}

	t.Run(fmt.Sprintf("Chirps Get Handler Test Case: No Chirps in Database"), func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/chirps/", nil)

		cfg.getChirpHandler(w, r)

		resp, _ := io.ReadAll(w.Body)

		if got := string(resp); got != empty.wantBody {
			t.Fatalf("Test failed (empty chirp list): got %q, want %q", got, empty.wantBody)
		}

		if w.Code != empty.wantCode {
			t.Fatalf("Test failed (no content code): got %d, want %d", w.Code, empty.wantCode)
		}
	})

	postW := httptest.NewRecorder()
	firstChirp := httptest.NewRequest("POST", "/api/chirps", strings.NewReader(`{"body": "My first Chirp"}`))
	secondChirp := httptest.NewRequest("POST", "/api/chirps", strings.NewReader(`{"body": "My second Chirp"}`))
	cfg.postChirpHandler(postW, firstChirp)
	cfg.postChirpHandler(postW, secondChirp)

	cases := []struct {
		wantCode int
		wantBody string
	}{
		{
			wantCode: 200,
			wantBody: `[{"id": 1, "body": "My first Chirp"}, {"id": 2, "body": "My second Chirp"}]`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Chirps Get Handler Test Case %d", i), func(t *testing.T) {

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/chirps/", nil)

			cfg.getChirpHandler(w, r)

			resp, _ := io.ReadAll(w.Body)

			if got := string(resp); got != c.wantBody {

			}
		})
	}

	err = os.Remove(testDbName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not cleanup database file: %q", err)
	}
}
