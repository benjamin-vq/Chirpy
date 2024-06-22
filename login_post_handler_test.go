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

func TestLoginPostHandler(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := database.NewDB(testDbName)
	if err != nil {
		log.Fatalf("Could not create new database for test: %q", err)
	}

	cfg := apiConfig{
		DB: db,
	}

	users := []struct {
		request string
	}{
		{
			request: `{"email": "loginhandler@chirpy.com", "password": "hello123"}`,
		},
		{
			request: `{"email": "2nduser@handler.io", "password": "pass"}`,
		},
	}

	for _, user := range users {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/users", strings.NewReader(user.request))
		cfg.postUsersHandler(w, req)
	}

	cases := []struct {
		wantCode int
		request  string
		wantBody string
	}{
		{
			wantCode: 200,
			request:  `{"email": "loginhandler@chirpy.com", "password": "hello123"}`,
			wantBody: `{"email":"loginhandler@chirpy.com","id":1}`,
		},
		{
			wantCode: 200,
			request:  `{"email": "2nduser@handler.io", "password": "pass"}`,
			wantBody: `{"email":"2nduser@handler.io","id":2}`,
		},
		{
			wantCode: 401,
			request:  `{"email": "loginhandler@chirpy.com", "password": "wrongpassword"}`,
			wantBody: `{"error":"Unauthorized"}`,
		},
		{
			wantCode: 401,
			request:  `{"email": "userNotExists@chirpy.com", "password": "any"}`,
			wantBody: `{"error":"Unauthorized"}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Login Post Handler Test Case %d", i), func(t *testing.T) {

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/login", strings.NewReader(c.request))

			cfg.loginPostHandler(w, req)

			resp, _ := io.ReadAll(w.Body)
			if got := string(resp); got != c.wantBody {
				t.Errorf("Test failed (email): got %q, want %q", got, c.wantBody)
			}
			if got := w.Code; got != c.wantCode {
				t.Errorf("Test failed (code): got %d, want %d", got, c.wantCode)
			}
		})
	}

	err = os.Remove(testDbName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not cleanup database file: %q", err)
	}
}
