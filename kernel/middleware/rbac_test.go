package middleware

import (
	"context"
	"errors"
	"labflow/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubChecker struct {
	allowed bool
	err     error
}

func (s stubChecker) HasAccess(context.Context, int, string) (bool, error) { return s.allowed, s.err }

func TestRBACMiddlewareStatus(t *testing.T) {
	cases := []struct {
		name    string
		checker stubChecker
		want    int
	}{
		{"allowed", stubChecker{allowed: true}, http.StatusOK},
		{"denied", stubChecker{allowed: false}, http.StatusForbidden},
		{"check failed", stubChecker{err: errors.New("db down")}, http.StatusServiceUnavailable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := httpx.Chain(okHandler, asUser(1), RBACMiddleware(tc.checker, "read"))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
			if w.Code != tc.want {
				t.Fatalf("got %d, want %d", w.Code, tc.want)
			}
		})
	}
}

func TestRBACMiddlewareWithoutUser(t *testing.T) {
	h := RBACMiddleware(stubChecker{allowed: true}, "read")(okHandler)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
