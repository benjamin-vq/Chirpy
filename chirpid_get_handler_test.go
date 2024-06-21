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

func TestGetChirpIdHandler(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cases := []struct {
		code int
		id   string
		want string
	}{
		{
			code: 200,
			id:   "1",
			want: `{"body":"A good chirp","id":1}`,
		},
		{
			code: 400,
			id:   "invalid",
			want: `{"error":"Provided id is not valid"}`,
		},
		{
			code: 404,
			id:   "27",
			want: `{"error":"chirp with id 27 does not exist"}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Chirp by Id Handler Test Case %d", i), func(t *testing.T) {
			db, err := database.NewDB(testDbName)
			if err != nil {
				log.Fatalf("Could not create new database for test case %d: %q", i, err)
			}

			cfg := apiConfig{
				DB: db,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "http://chirpy.com", strings.NewReader(`{"body":"A good chirp"}`))
			cfg.postChirpHandler(w, req)

			idW := httptest.NewRecorder()
			idReq := httptest.NewRequest("GET", "/api/chirps/", nil)
			idReq.SetPathValue("chirpId", c.id)

			cfg.chirpIdGetHandler(idW, idReq)

			resp, _ := io.ReadAll(idW.Body)

			if got := string(resp); got != c.want {
				t.Errorf("Test failed (id): got %q, want %q", got, c.want)
			}

			if got := idW.Code; got != c.code {
				t.Errorf("Test failed (status code): got %d, want %d", got, c.code)

			}

			err = os.Remove(testDbName)
			if err != nil {
				t.Fatalf("Could not delete test database for next test: %q", err.Error())
			}
		})
	}
}
