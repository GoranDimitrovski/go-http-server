package http

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"simplesurance/internal/application"
	"simplesurance/internal/infrastructure/repository"
)

func TestRoutes(t *testing.T) {
	store := repository.NewMemoryStore(filepath.Join(t.TempDir(), "ts.log"))
	h := NewTimestampHandler(application.NewTimestampService(store, 60), log.New(io.Discard, "", 0))
	mux := http.NewServeMux()
	h.Routes(mux, "/")

	for _, tc := range []struct {
		method, path string
		status       int
		body         string
	}{
		{"GET", "/", 200, `{"count":1}` + "\n"},
		{"GET", "/", 200, `{"count":2}` + "\n"},
		{"POST", "/", 405, ""},
		{"GET", "/health", 200, "OK"},
	} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.status || (tc.body != "" && rec.Body.String() != tc.body) {
			t.Errorf("%s %s = %d %q; want %d %q", tc.method, tc.path, rec.Code, rec.Body.String(), tc.status, tc.body)
		}
	}
}
