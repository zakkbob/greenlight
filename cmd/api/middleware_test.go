package main

import (
	"net/http"
	"testing"

	"github.com/zakkbob/greenlight/internal/assert"
)

func TestRateLimit(t *testing.T) {
	cfg := newTestConfig(t)
	app := newTestApplication(t, cfg)
	ts := newTestServer(t, app.routes())

	for range 5 {
		res := ts.get(t, "/v1/healthcheck")
		assert.Equal(t, res.status, http.StatusOK)
	}

	res := ts.get(t, "/v1/healthcheck")
	assert.Equal(t, res.status, http.StatusTooManyRequests)
}
