package main

import (
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/benjamin-vq/chirpy/internal/database"
)

const testDbName = "testDatabase.json"

func TestPostChirpHandler(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cases := []struct {
		code int
		body string
		want string
	}{
		{
			code: 201,
			body: `{"body": "A good chirp"}`,
			want: `{"body":"A good chirp","id":1}`,
		},
		{
			code: 201,
			body: `{"body": "A decent chirp, chirped by fornax (not Fornax)"}`,
			want: `{"body":"A decent chirp, chirped by **** (not ****)","id":1}`,
		},
		{
			code: 400,
			body: `{"body": "A really really long, omnipotent chirp, one may call it the best chirp. Capable of surpassing the longest of limits, beyond human imagination."}`,
			want: `{"error":"chirp length exceeds limit"}`,
		},
		{
			code: 500,
			body: `invalid json`,
			want: `{"error":"Could not decode chirp"}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Chirps Post Handler Test Case %d", i), func(t *testing.T) {

			db, err := database.NewDB(testDbName)
			if err != nil {
				log.Fatalf("Could not create new database for test case %d: %q", i, err)
			}

			cfg := apiConfig{
				DB: db,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "http://chirpy.com", strings.NewReader(c.body))

			cfg.postChirpHandler(w, req)

			resp, _ := io.ReadAll(w.Body)
			if got := string(resp); got != c.want {
				t.Errorf("Test failed (body): got %q, want %q", got, c.want)
			}
			if got := w.Code; got != c.code {
				t.Errorf("Test failed (code): got %d, want %d", got, c.code)
			}

			err = os.Remove(testDbName)
			if err != nil {
				t.Fatalf("Could not delete test database for next test: %q", err.Error())
			}
		})
	}
}

func TestReplaceBadWords(t *testing.T) {

	cases := []struct {
		body string
		want string
	}{
		{
			body: "No bad words here!",
			want: "No bad words here!",
		}, {
			body: "A kerfuffle sounds very good",
			want: "A **** sounds very good",
		}, {
			body: "What is a sharbert?",
			want: "What is a ****?",
		}, {
			body: "A new pokemon was announced: fornax",
			want: "A new pokemon was announced: ****",
		}, {
			body: "A Fornax was caught eating a Kerfuffle in a Sharbert",
			want: "A **** was caught eating a **** in a ****",
		}, {
			body: "I really need a kerfuffle to go to bed sooner, Fornax !",
			want: "I really need a **** to go to bed sooner, **** !",
		}, {
			body: "",
			want: "",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {

			sanitized := replaceBadWords(c.body)

			if sanitized != c.want {
				t.Errorf("Test failed: got %q but want %q", sanitized, c.want)
				return
			}
		})
	}
}
