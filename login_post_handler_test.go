package main

import (
	"encoding/json"
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

	t.Setenv("JWT_SECRET", "dGVzdA==")

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
		wantMap  map[string]any
	}{
		{
			wantCode: 200,
			request:  `{"email": "loginhandler@chirpy.com", "password": "hello123"}`,
			wantMap: map[string]any{
				"email": "loginhandler@chirpy.com",
				"id":    float64(1),
			},
		},
		{
			wantCode: 200,
			request:  `{"email": "2nduser@handler.io", "password": "pass"}`,
			wantMap: map[string]any{
				"email": "2nduser@handler.io",
				"id":    float64(2),
			},
		},
		{
			wantCode: 401,
			request:  `{"email": "loginhandler@chirpy.com", "password": "wrongpassword"}`,
			wantMap: map[string]any{
				"error": "Unauthorized",
			},
		},
		{
			wantCode: 401,
			request:  `{"email": "userNotExists@chirpy.com", "password": "any"}`,
			wantMap: map[string]any{
				"error": "Unauthorized",
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Login Post Handler Test Case %d", i), func(t *testing.T) {

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/login", strings.NewReader(c.request))

			cfg.loginPostHandler(w, req)

			resp, _ := io.ReadAll(w.Body)

			var jsonResponse map[string]any
			if err := json.Unmarshal(resp, &jsonResponse); err != nil {
				t.Fatalf("Could not unmarshal response: %q", err)
			}

			if _, ok := jsonResponse["token"]; !ok && w.Code == 200 {
				t.Errorf("Test failed, expected response to have a 'token' field")
			}
			delete(jsonResponse, "token")

			for k, want := range c.wantMap {
				if actual, _ := jsonResponse[k]; actual != want {
					t.Errorf("Test failed, comparing json response got %v, want %v", actual, want)
				}
			}

		})
	}

	err = os.Remove(testDbName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not cleanup database file: %q", err)
	}
}
