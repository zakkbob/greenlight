package main

import (
	"net/http"
	"testing"

	"github.com/zakkbob/greenlight/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	cfg := newTestConfig(t)
	app := newTestApplication(t, cfg)
	ts := newTestServer(t, app.routes())
	defer ts.close()

	var js struct {
		Status     string `json:"status"`
		SystemInfo struct {
			Environment string `json:"environment"`
			Version     string `json:"version"`
		} `json:"system_info"`
	}

	res := ts.getJSON(t, "/v1/healthcheck", &js)

	expectedVary := []string{"Access-Control-Request-Method", "Authorization", "Origin"}
	assert.EqualSlicesUnordered(t, res.header.Values("Vary"), expectedVary)

	assert.Equal(t, res.header.Get("Content-Type"), "application/json")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, js.Status, "available")
	assert.Equal(t, js.SystemInfo.Version, version)
}
