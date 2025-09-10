package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/zakkbob/greenlight/internal/assert"
	"github.com/zakkbob/greenlight/internal/data"
)

func TestRateLimit(t *testing.T) {
	app := newTestApplication(t)
	app.config = newTestConfig(t)
	ts := newTestServer(t, app.routes())

	for range 5 {
		res := ts.get(t, "/v1/healthcheck", nil)
		assert.Equal(t, res.status, http.StatusOK)
	}

	res := ts.get(t, "/v1/healthcheck", nil)
	assert.Equal(t, res.status, http.StatusTooManyRequests)
}

func TestRegisterUser(t *testing.T) {
	if testing.Short() {
		t.Skip("middleware: skipping end-to-end integration test")
	}

	app := newTestApplication(t)
	app.models = newTestModels(t, newTestDB(t))
	app.config = newTestConfig(t)
	ts := newTestServer(t, app.routes())

	body := map[string]string{
		"name":     "john",
		"email":    "john@example.com",
		"password": "pa55word",
	}

	var js struct {
		User data.User `json:"user"`
	}

	res := ts.post(t, "/v1/users", nil, body)
	res.Decode(t, &js)

	assert.Equal(t, js.User.Name, "john")
	assert.Equal(t, js.User.Email, "john@example.com")
	assert.Equal(t, js.User.ID, 1)
	assert.Equal(t, js.User.Activated, false)
}

func TestUserPermissions(t *testing.T) {
	tests := []struct {
		name          string
		permissions   data.Permissions
		activated     bool
		authenticated bool
		path          string
		status        int
	}{
		{
			name:          "No auth required or given",
			path:          "/v1/healthcheck",
			activated:     false,
			authenticated: false,
			permissions:   data.Permissions{},
			status:        http.StatusOK,
		},
		{
			name:          "Auth required but not given",
			path:          "/v1/movies",
			activated:     false,
			authenticated: false,
			permissions:   data.Permissions{},
			status:        http.StatusUnauthorized,
		},
		{
			name:          "Doesn't have permission",
			path:          "/v1/movies",
			activated:     true,
			authenticated: true,
			permissions:   data.Permissions{},
			status:        http.StatusForbidden,
		},
		{
			name:          "Has permission",
			path:          "/v1/movies",
			activated:     true,
			authenticated: true,
			permissions:   data.Permissions{"movies:read"},
			status:        http.StatusOK,
		},
		{
			name:          "Wrong permission",
			path:          "/v1/movies",
			activated:     true,
			authenticated: true,
			permissions:   data.Permissions{"debug:read"},
			status:        http.StatusForbidden,
		},
		{
			name:          "Account not activated",
			path:          "/v1/movies",
			activated:     false,
			authenticated: true,
			permissions:   data.Permissions{"movies:read"},
			status:        http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t)
			app.config = newTestConfig(t)
			app.models = newTestModels(t, newTestDB(t))
			ts := newTestServer(t, app.routes())

			header := http.Header{}

			if tt.authenticated {
				user := &data.User{
					Name:      "John",
					Email:     "john@example.com",
					Activated: tt.activated,
				}
				user.Password.Set("pa55word")
				err := app.models.Users.Insert(user)
				if err != nil {
					t.Fatal(err)
				}

				err = app.models.Permissions.AddForUser(user.ID, tt.permissions...)
				if err != nil {
					t.Fatal(err)
				}

				token, err := app.models.Tokens.New(user.ID, time.Hour, data.ScopeAuthentication)
				if err != nil {
					t.Fatal(err)
				}

				header.Set("Authorization", "Bearer "+token.Plaintext)
			}

			res := ts.get(t, tt.path, header)
			t.Log(res.body)
			assert.Equal(t, res.status, tt.status)
		})
	}
}
