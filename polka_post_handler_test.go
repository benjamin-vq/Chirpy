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

func TestPolkaPostHandler(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := database.NewDB(testDbName)
	if err != nil {
		log.Fatalf("Could not create new database for test: %q", err)
	}

	cfg := apiConfig{
		DB:          db,
		polkaApiKey: "1234",
	}

	user := `{"email": "newuser@chirpy.com", "password": "hey!"}`
	createW := httptest.NewRecorder()
	createReq := httptest.NewRequest("POST", "/api/users", strings.NewReader(user))
	cfg.postUsersHandler(createW, createReq)

	cases := []struct {
		request  string
		wantCode int
		wantBody string
	}{
		{
			request:  `{"event": "user.upgraded","data": {"user_id": 1}}`,
			wantCode: 204,
			wantBody: `""`,
		},
		{
			request:  `{"event": "another.event","data": {"user_id": 1}}`,
			wantCode: 204,
			wantBody: `""`,
		},
		{
			request:  `{"event": "another.event","data": {"user_id": 2}}`,
			wantCode: 204,
			wantBody: `""`,
		},
		{
			request:  `{"event": "user.upgraded","data": {"user_id": 2}}`,
			wantCode: 404,
			wantBody: `{}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Polka Post Handler Test Case %d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/polka/webhooks", strings.NewReader(c.request))
			req.Header.Set("Authorization", "ApiKey "+cfg.polkaApiKey)

			cfg.postPolkaHandler(w, req)

			resp, _ := io.ReadAll(w.Body)

			if got := string(resp); got != c.wantBody {
				t.Errorf("Test case failed (body): got %s, want %s", got, c.wantBody)
			}

			if w.Code != c.wantCode {
				t.Errorf("Test case failed (body): got %d, want %d", w.Code, c.wantCode)
			}
		})
	}

	err = os.Remove(testDbName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Could not cleanup database file: %q", err)
	}
}
