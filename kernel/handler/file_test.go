package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"labflow/service"

	"github.com/gin-gonic/gin"
)

type fakeFetcher struct {
	data      *service.ObjectData
	err       error
	gotBucket string
	gotObject string
	called    bool
}

func (f *fakeFetcher) GetObject(ctx context.Context, bucket, object string) (*service.ObjectData, error) {
	f.called = true
	f.gotBucket = bucket
	f.gotObject = object
	return f.data, f.err
}

func newTestRouter(h *FileHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/files/reservations/:object", h.ServeReservationFile)
	r.GET("/api/files/users/:object", h.ServeUserFile)
	r.GET("/api/files/certificates/:object", h.ServeCertificateFile)
	return r
}

func TestServeReservationFileOK(t *testing.T) {
	f := &fakeFetcher{data: &service.ObjectData{
		Body:        io.NopCloser(strings.NewReader("PNGBYTES")),
		ContentType: "image/png",
		Size:        8,
	}}
	h := NewFileHandler(f, "reservations", "users", "certificates")
	r := newTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/files/reservations/abc.png", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "PNGBYTES" {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "image/png" {
		t.Fatalf("unexpected content-type: %q", ct)
	}
	if f.gotBucket != "reservations" || f.gotObject != "abc.png" {
		t.Fatalf("fetcher got bucket=%q object=%q", f.gotBucket, f.gotObject)
	}
}

func TestServeUserFileUsesUserBucket(t *testing.T) {
	f := &fakeFetcher{data: &service.ObjectData{
		Body:        io.NopCloser(strings.NewReader("JPG")),
		ContentType: "image/jpeg",
		Size:        3,
	}}
	h := NewFileHandler(f, "reservations", "users", "certificates")
	r := newTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/files/users/avatar.jpg", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if f.gotBucket != "users" {
		t.Fatalf("expected users bucket, got %q", f.gotBucket)
	}
}

func TestServeReservationFileNotFound(t *testing.T) {
	f := &fakeFetcher{err: service.ErrObjectNotFound}
	h := NewFileHandler(f, "reservations", "users", "certificates")
	r := newTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/files/reservations/missing.png", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestServeReservationFileInternalError(t *testing.T) {
	f := &fakeFetcher{err: errors.New("boom")}
	h := NewFileHandler(f, "reservations", "users", "certificates")
	r := newTestRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/files/reservations/x.png", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestValidObjectName(t *testing.T) {
	cases := map[string]bool{
		"abc.png":     true,
		"e5a3-1b.jpg": true,
		"":            false,
		"a/b.png":     false,
		"..":          false,
		"a..b":        false,
		"x\\y":        false,
	}
	for name, want := range cases {
		if got := validObjectName(name); got != want {
			t.Errorf("validObjectName(%q) = %v, want %v", name, got, want)
		}
	}
}
