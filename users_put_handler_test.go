package main

import (
	"encoding/json"
	"errors"
	"github.com/benjamin-vq/chirpy/internal/database"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPutUsersHandler(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := database.NewDB(testDbName)
	if err != nil {
		log.Fatalf("Could not create new database for test: %q", err)
	}

	cfg := apiConfig{
		DB:        db,
		jwtSecret: "dGVzdA==",
	}

	user := `{"email": "newuser@chirpy.com", "password": "hey!"}`

	createW := httptest.NewRecorder()
	createReq := httptest.NewRequest("POST", "/api/users", strings.NewReader(user))
	cfg.postUsersHandler(createW, createReq)

	loginRequest := `{"email": "newuser@chirpy.com", "password": "hey!", "expires_in_seconds": 5}`
	loginW := httptest.NewRecorder()
	loginReq := httptest.NewRequest("POST", "/api/login", strings.NewReader(loginRequest))
	cfg.loginPostHandler(loginW, loginReq)

	loginResp := map[string]string{}
	decoder := json.NewDecoder(loginW.Body)
	err = decoder.Decode(&loginResp)

	token, _ := loginResp["token"]
	want := `{"email":"updated@user.com","id":1}`

	putW := httptest.NewRecorder()
	putReq := httptest.NewRequest("PUT", "/api/users", strings.NewReader(want))
	putReq.Header.Set("Authorization", "Bearer "+token)

	cfg.putUsersHandler(putW, putReq)

	t.Run("Updated User Test", func(t *testing.T) {

		resp, _ := io.ReadAll(putW.Body)

		if string(resp) != want {
			t.Fatalf("Incorrect update response: got %s, want %s", string(resp), want)
		}
	})

	err = os.Remove(testDbName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not cleanup database file: %q", err)
	}
}
