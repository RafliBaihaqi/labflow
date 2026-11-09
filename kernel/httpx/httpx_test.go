package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSON(w, http.StatusCreated, map[string]string{"msg": "<ok>"})

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content type = %q", ct)
	}
	// Same bytes as json.Marshal (HTML escaped, no trailing newline), as Gin wrote.
	if got, want := w.Body.String(), `{"msg":"\u003cok\u003e"}`; got != want {
		t.Fatalf("body = %q", got)
	}
}

func TestChainOrder(t *testing.T) {
	var order []string
	mw := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	h := Chain(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { order = append(order, "handler") }),
		mw("first"), mw("second"))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if got := strings.Join(order, ","); got != "first,second,handler" {
		t.Fatalf("order = %s", got)
	}
}

func TestUserID(t *testing.T) {
	if _, ok := UserID(context.Background()); ok {
		t.Fatal("empty context must not carry a user")
	}
	id, ok := UserID(WithUserID(context.Background(), 42))
	if !ok || id != 42 {
		t.Fatalf("got %d, %v", id, ok)
	}
}

func TestQueryDefault(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?page=&limit=5", nil)
	if got := QueryDefault(r, "page", "1"); got != "" {
		t.Errorf("present but empty: got %q", got)
	}
	if got := QueryDefault(r, "limit", "10"); got != "5" {
		t.Errorf("present: got %q", got)
	}
	if got := QueryDefault(r, "sort", "asc"); got != "asc" {
		t.Errorf("absent: got %q", got)
	}
}

func TestClientIP(t *testing.T) {
	cases := []struct {
		name   string
		remote string
		xff    string
		realIP string
		want   string
	}{
		{"direct client", "203.0.113.7:5000", "", "", "203.0.113.7"},
		{"public peer cannot spoof", "203.0.113.7:5000", "1.2.3.4", "", "203.0.113.7"},
		{"behind local proxy", "127.0.0.1:5000", "198.51.100.9", "", "198.51.100.9"},
		{"proxy chain", "172.18.0.1:5000", "198.51.100.9, 10.0.0.5", "", "198.51.100.9"},
		{"client-supplied hop is skipped", "172.18.0.1:5000", "1.2.3.4, 198.51.100.9", "", "198.51.100.9"},
		{"x-real-ip", "127.0.0.1:5000", "", "198.51.100.9", "198.51.100.9"},
		{"local peer without headers", "127.0.0.1:5000", "", "", "127.0.0.1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.RemoteAddr = tc.remote
			if tc.xff != "" {
				r.Header.Set("X-Forwarded-For", tc.xff)
			}
			if tc.realIP != "" {
				r.Header.Set("X-Real-IP", tc.realIP)
			}
			if got := ClientIP(r); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestRoute(t *testing.T) {
	mux := http.NewServeMux()
	var got string
	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) { got = Route(r) })
	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/items/9", nil))
	if got != "/items/{id}" {
		t.Fatalf("route = %q", got)
	}
}

func TestStatusWriter(t *testing.T) {
	sw := NewStatusWriter(httptest.NewRecorder())
	if NewStatusWriter(sw) != sw {
		t.Fatal("wrapping twice must reuse the writer")
	}
	if sw.Written() || sw.Status() != http.StatusOK {
		t.Fatal("unwritten writer reports 200 and Written() == false")
	}
	_, _ = sw.Write([]byte("abc"))
	if !sw.Written() || sw.Status() != http.StatusOK || sw.Size() != 3 {
		t.Fatalf("after write: written=%v status=%d size=%d", sw.Written(), sw.Status(), sw.Size())
	}
}
