package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsLabelsByRouteTemplate(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /items/{id}", okHandler)
	h := Metrics(mux)

	before := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/items/{id}", "200"))
	for _, id := range []string{"1", "2", "3"} {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/items/"+id, nil))
	}
	after := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/items/{id}", "200"))

	if after-before != 3 {
		t.Fatalf("expected 3 requests under one route label, got %v", after-before)
	}
}

func TestMetricsUnmatchedRoute(t *testing.T) {
	h := Metrics(http.NewServeMux())

	before := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "unmatched", "404"))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/nope", nil))
	after := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "unmatched", "404"))

	if after-before != 1 {
		t.Fatalf("expected unmatched 404 to be counted once, got %v", after-before)
	}
}

// The route label must survive middleware that sits between Metrics and the mux.
func TestMetricsRouteThroughOtherMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /things/{id}", okHandler)
	h := Metrics(CORS(CORSConfig{})(mux))

	before := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/things/{id}", "200"))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/things/7", nil))
	after := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/things/{id}", "200"))

	if after-before != 1 {
		t.Fatalf("expected the route template label, got delta %v", after-before)
	}
}
