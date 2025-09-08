package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/zakkbob/greenlight/internal/data"
)

type testResponse struct {
	status int
	body   string
	header http.Header
}

func (res *testResponse) Decode(t *testing.T, v any) {
	d := json.NewDecoder(bytes.NewBuffer([]byte(res.body)))
	d.DisallowUnknownFields()
	err := d.Decode(v)
	if err != nil {
		t.Fatalf("Failed to decode JSON '%v': %v", res.body, err)
	}
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

func newTestModels(t *testing.T, db *sql.DB) data.Models {
	return data.Models{
		Movies:      data.MovieModel{DB: db},
		Permissions: data.PermissionModel{DB: db},
		Tokens:      data.TokenModel{DB: db},
		Users:       data.UserModel{DB: db},
	}
}

func newTestApplication(t *testing.T) application {
	return application{
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

func processResponse(t *testing.T, res *http.Response) testResponse {
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

func (ts *testServer) get(t *testing.T, path string) testResponse {
	res, err := http.Get(ts.server.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	return processResponse(t, res)
}

func (ts *testServer) post(t *testing.T, path string, data any) testResponse {
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to encode JSON data '%v': %v", data, err)
	}

	res, err := http.Post(ts.server.URL+path, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	return processResponse(t, res)
}

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://test_greenlight:password@localhost:5432/test_greenlight?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres", driver)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := m.Drop()
		if err != nil {
			t.Fatal(err)
		}

		srcErr, dbErr := m.Close()
		if srcErr != nil {
			t.Fatal(srcErr)
		}
		if dbErr != nil {
			t.Fatal(dbErr)
		}
	})

	version, _, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		t.Fatal(err)
	}

	if version != 0 || !errors.Is(err, migrate.ErrNilVersion) {
		m.Down()
		if err != nil {
			t.Fatal(err)
		}
	}

	err = m.Up()
	if err != nil {
		t.Fatal(err)
	}

	return db
}
