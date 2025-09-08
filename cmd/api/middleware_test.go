package main

import (
	"net/http"
	"testing"

	"github.com/zakkbob/greenlight/internal/assert"
	"github.com/zakkbob/greenlight/internal/data"
)

func TestRateLimit(t *testing.T) {
	app := newTestApplication(t)
	app.config = newTestConfig(t)
	ts := newTestServer(t, app.routes())

	for range 5 {
		res := ts.get(t, "/v1/healthcheck")
		assert.Equal(t, res.status, http.StatusOK)
	}

	res := ts.get(t, "/v1/healthcheck")
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

	res := ts.post(t, "/v1/users", body)
	res.Decode(t, &js)

	assert.Equal(t, js.User.Name, "john")
	assert.Equal(t, js.User.Email, "john@example.com")
	assert.Equal(t, js.User.ID, 1)
	assert.Equal(t, js.User.Activated, false)
}
