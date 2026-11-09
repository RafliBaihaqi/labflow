package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewHandler registers every route on one ServeMux, which panics on
// conflicting patterns. Building it is the test; the requests check that the
// global middleware and the method-aware routing behave as the API expects.
func TestNewHandlerRoutes(t *testing.T) {
	t.Setenv("DISABLE_RATE_LIMIT", "true")
	h := NewHandler("test-cookie-key-at-least-32-bytes", nil, nil, nil, nil)

	cases := []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/health", http.StatusOK},
		{http.MethodHead, "/health", http.StatusOK},
		{http.MethodGet, "/api/reservations", http.StatusUnauthorized},
		{http.MethodGet, "/api/reservations/stats", http.StatusUnauthorized},
		{http.MethodGet, "/api/auth/me", http.StatusUnauthorized},
		{http.MethodDelete, "/api/reservations", http.StatusMethodNotAllowed},
		{http.MethodGet, "/nope", http.StatusNotFound},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.want {
			t.Errorf("%s %s = %d, want %d", tc.method, tc.path, w.Code, tc.want)
		}
	}
}
