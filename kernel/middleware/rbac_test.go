package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubChecker struct {
	allowed bool
	err     error
}

func (s stubChecker) HasAccess(context.Context, int, string) (bool, error) { return s.allowed, s.err }

func TestRBACMiddlewareStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
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
			r := gin.New()
			r.GET("/x", func(c *gin.Context) { c.Set("id", 1); c.Next() }, RBACMiddleware(tc.checker, "read"), func(c *gin.Context) { c.Status(http.StatusOK) })
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
			if w.Code != tc.want {
				t.Fatalf("got %d, want %d", w.Code, tc.want)
			}
		})
	}
}
