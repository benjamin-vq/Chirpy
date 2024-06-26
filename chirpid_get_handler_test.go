package main

import (
	"encoding/json"
	"errors"
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

	db, err := database.NewDB(testDbName)
	if err != nil {
		log.Fatalf("Could not create new database for test: %q", err)
	}

	cfg := apiConfig{
		DB: db,
	}

	user := `{"email": "newuser@chirpy.com", "password": "hey!"}`
	createW := httptest.NewRecorder()
	createReq := httptest.NewRequest("POST", "/api/users", strings.NewReader(user))
	cfg.postUsersHandler(createW, createReq)

	loginRequest := `{"email": "newuser@chirpy.com", "password": "hey!"}`
	loginW := httptest.NewRecorder()
	loginReq := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginRequest))
	cfg.loginPostHandler(loginW, loginReq)

	loginResp := map[string]string{}
	decoder := json.NewDecoder(loginW.Body)
	err = decoder.Decode(&loginResp)

	token, _ := loginResp["token"]

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://chirpy.com", strings.NewReader(`{"body":"A good chirp"}`))
	req.Header.Add("Authorization", "Bearer "+token)
	cfg.postChirpHandler(w, req)

	cases := []struct {
		code int
		id   string
		want string
	}{
		{
			code: 200,
			id:   "1",
			want: `{"body":"A good chirp","id":1,"author_id":1}`,
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
		})
	}

	err = os.Remove(testDbName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not cleanup database file: %q", err)
	}
}
