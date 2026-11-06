package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsLabelsByRouteTemplate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Metrics())
	r.GET("/items/:id", func(c *gin.Context) { c.Status(http.StatusOK) })

	before := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/items/:id", "200"))
	for _, id := range []string{"1", "2", "3"} {
		r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/items/"+id, nil))
	}
	after := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/items/:id", "200"))

	if after-before != 3 {
		t.Fatalf("expected 3 requests under one route label, got %v", after-before)
	}
}

func TestMetricsUnmatchedRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Metrics())

	before := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "unmatched", "404"))
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/nope", nil))
	after := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "unmatched", "404"))

	if after-before != 1 {
		t.Fatalf("expected unmatched 404 to be counted once, got %v", after-before)
	}
}
