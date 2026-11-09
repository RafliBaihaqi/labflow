package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testCORS() http.Handler {
	return CORS(CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})(okHandler)
}

func TestCORSPreflightFromAllowedOrigin(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "http://api.local/api/reservations", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	testCORS().ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
	want := map[string]string{
		"Access-Control-Allow-Origin":      "http://localhost:5173",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "GET,POST",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Max-Age":           "43200",
	}
	for k, v := range want {
		if got := w.Header().Get(k); got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}
	if w.Body.Len() != 0 {
		t.Errorf("preflight must not reach the handler, body = %q", w.Body.String())
	}
}

func TestCORSSimpleRequestFromAllowedOrigin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://api.local/api/reservations", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	testCORS().ServeHTTP(w, req)

	if w.Code != http.StatusOK || w.Body.String() != "OK" {
		t.Fatalf("got %d %q, want 200 OK", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Access-Control-Expose-Headers"); got != "Content-Length" {
		t.Errorf("Access-Control-Expose-Headers = %q", got)
	}
	if got := w.Header().Get("Vary"); got != "Origin" {
		t.Errorf("Vary = %q, want Origin", got)
	}
}

func TestCORSRejectsUnknownOrigin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://api.local/api/reservations", nil)
	req.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	testCORS().ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestCORSIgnoresSameOriginAndNoOrigin(t *testing.T) {
	for _, origin := range []string{"", "http://api.local"} {
		req := httptest.NewRequest(http.MethodPost, "http://api.local/auth/login", nil)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		w := httptest.NewRecorder()
		testCORS().ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("origin %q: status = %d, want 200", origin, w.Code)
		}
		if w.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("origin %q: unexpected CORS headers", origin)
		}
	}
}
