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

func TestPostUsersHandler(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cases := []struct {
		code int
		body string
		want string
	}{
		{
			code: 201,
			body: `{"email": "myemail@chirpy.com", "password": "test1234"}`,
			want: `{"email":"myemail@chirpy.com","id":1}`,
		},
		{
			code: 201,
			body: `{"email": "another@email.io", "id": 5958, "password": "1234567890"}`,
			want: `{"email":"another@email.io","id":1}`,
		},
		{
			code: 400,
			body: `{"email": "", "password": "$!\"#ña!"}`,
			want: `{"error":"email can not be empty"}`,
		},
		{
			code: 400,
			body: `{"email": "test@emails.org", "password": ""}`,
			want: `{"error":"password can not be empty"}`,
		},
		{
			code: 400,
			body: `{"email": "", "password": ""}`,
			want: `{"error":"email can not be empty"}`,
		},
		{
			code: 400,
			body: `{"email": "asdf"}`,
			want: `{"error":"password can not be empty"}`,
		},
		{
			code: 400,
			body: `{"password": "hey"}`,
			want: `{"error":"email can not be empty"}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Post Users Handler Test Case %d", i), func(t *testing.T) {

			db, err := database.NewDB(testDbName)
			if err != nil {
				log.Fatalf("Could not create new database for test case %d: %q", i, err)
			}

			cfg := apiConfig{
				DB: db,
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/users", strings.NewReader(c.body))

			cfg.postUsersHandler(w, req)

			resp, _ := io.ReadAll(w.Body)
			if got := string(resp); got != c.want {
				t.Errorf("Test failed (email): got %q, want %q", got, c.want)
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
