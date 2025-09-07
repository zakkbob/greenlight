package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type testResponse struct {
	status int
	body   string
	header http.Header
}

func newTestConfig(t *testing.T) config {
	return config{
		port: 4000,
		env:  "production",
		limiter: struct {
			rps     float64
			burst   int
			enabled bool
		}{
			rps:     1,
			burst:   5,
			enabled: true,
		},
		metrics: false,
	}
}

func newTestApplication(t *testing.T, cfg config) application {
	return application{
		config: cfg,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		wg:     sync.WaitGroup{},
	}
}

type testServer struct {
	server *httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) testServer {
	return testServer{
		server: httptest.NewServer(h),
	}
}

func (ts *testServer) close() {
	ts.server.Close()
}

func (ts *testServer) get(t *testing.T, path string) testResponse {
	t.Helper()
	res, err := http.Get(ts.server.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status: res.StatusCode,
		body:   string(body),
		header: res.Header,
	}

}

func (ts *testServer) getJSON(t *testing.T, path string, v any) testResponse {
	t.Helper()
	res, err := http.Get(ts.server.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	d := json.NewDecoder(bytes.NewBuffer(body))
	d.DisallowUnknownFields()
	err = d.Decode(v)
	if err != nil {
		t.Fatalf("Failed to decode JSON '%v': %v", string(body), err)
	}

	return testResponse{
		status: res.StatusCode,
		body:   string(body),
		header: res.Header,
	}

}
