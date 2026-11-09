// Package microbench measures only the HTTP layer: routing, middleware and JSON
// encoding, with no database. It compares the Gin setup LabFlow used with the
// net/http setup that replaced it, serving the same response.
//
//	cd loadtest/microbench && go test -bench . -benchmem -count 6
package microbench

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"labflow/dto"
	"labflow/httpx"
	"labflow/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var origins = []string{"http://localhost:5173", "http://localhost:3000"}

func sampleReservation() dto.ReservationResponse {
	note, status := "keep cold", "pending"
	now := time.Date(2026, 11, 9, 10, 0, 0, 0, time.UTC)
	return dto.ReservationResponse{
		ID: 42, Xid: uuid.MustParse("41b3a272-0775-413e-94ab-e838c6347466"), UserID: 4,
		ReservationDate: now, Status: &status, SampleName: "Purified water", SampleAmount: 3,
		SpecialInstructions: &note, SampleType: "Liquid", SampleTakenDate: now, SampleTakenBy: "k6",
		CreatedAt: &now, UpdatedAt: &now, Message: "Reservation fetched",
		Tests: []dto.TestResponse{
			{ID: 1, Method: "USP <61>", Criteria: "< 100 CFU/g", TestUnit: "CFU/g", TestList: "Total Aerobic Microbial Count"},
			{ID: 2, Method: "USP <791>", Criteria: "6.5 - 7.5", TestUnit: "pH", TestList: "pH"},
		},
		User: dto.UserResponse{ID: 4, Email: "customer@labflow.local", Name: "Customer"},
	}
}

// ginRouter mirrors the old server: recovery, CORS, a session middleware that
// sets the user ID, then a handler that reads a path parameter and writes JSON.
func ginRouter() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	body := sampleReservation()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{AllowOrigins: origins, AllowCredentials: true, MaxAge: 12 * time.Hour}))
	auth := func(c *gin.Context) { c.Set("id", 4); c.Next() }
	api := r.Group("/api", auth)
	api.GET("/reservations/:xid", func(c *gin.Context) {
		if _, ok := c.Get("id"); !ok || c.Param("xid") == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, body)
	})
	api.GET("/reservations/stats", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"total": 1}) })
	api.GET("/reservations", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"data": []int{}}) })
	return r
}

// stdRouter is the same pipeline on net/http with the API's own httpx helpers
// and middleware.
func stdRouter() http.Handler {
	body := sampleReservation()
	auth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(httpx.WithUserID(r.Context(), 4)))
		})
	}
	mux := http.NewServeMux()
	mux.Handle("GET /api/reservations/{xid}", auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := httpx.UserID(r.Context()); !ok || r.PathValue("xid") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, body)
	})))
	mux.Handle("GET /api/reservations/stats", auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]int{"total": 1})
	})))
	mux.Handle("GET /api/reservations", auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string][]int{"data": {}})
	})))
	return httpx.Chain(mux,
		middleware.Recover,
		middleware.CORS(middleware.CORSConfig{AllowOrigins: origins, AllowCredentials: true, MaxAge: 12 * time.Hour}),
	)
}

const (
	detailPath = "/api/reservations/41b3a272-0775-413e-94ab-e838c6347466" // ~1 KB JSON body
	statsPath  = "/api/reservations/stats"                                // 11-byte JSON body
)

func bench(b *testing.B, h http.Handler, path string) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Origin", "http://localhost:5173")
	b.ReportAllocs()
	for b.Loop() {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Fatalf("status %d", w.Code)
		}
	}
}

func BenchmarkGinDetail(b *testing.B)     { bench(b, ginRouter(), detailPath) }
func BenchmarkNetHTTPDetail(b *testing.B) { bench(b, stdRouter(), detailPath) }
func BenchmarkGinStats(b *testing.B)      { bench(b, ginRouter(), statsPath) }
func BenchmarkNetHTTPStats(b *testing.B)  { bench(b, stdRouter(), statsPath) }

// Both routers must return the same bytes, or the benchmark compares different work.
func TestSameResponse(t *testing.T) {
	get := func(h http.Handler) (int, http.Header, string) {
		req := httptest.NewRequest(http.MethodGet, "/api/reservations/41b3a272-0775-413e-94ab-e838c6347466", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		body, _ := io.ReadAll(w.Result().Body)
		return w.Code, w.Header(), string(body)
	}
	gc, gh, gb := get(ginRouter())
	sc, sh, sb := get(stdRouter())
	if gc != sc || gb != sb {
		t.Fatalf("responses differ:\ngin %d %s\nstd %d %s", gc, gb, sc, sb)
	}
	for _, k := range []string{"Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Vary"} {
		if gh.Get(k) != sh.Get(k) {
			t.Errorf("%s: gin %q, std %q", k, gh.Get(k), sh.Get(k))
		}
	}
}
