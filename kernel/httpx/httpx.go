// Package httpx holds the small helpers the API needs on top of net/http:
// JSON responses, middleware chaining, the authenticated user in the request
// context, client IP resolution and a status-recording ResponseWriter.
package httpx

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
)

// WriteJSON writes v as a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		log.Printf("[httpx] encode response: %v", err)
		http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// Middleware wraps a handler with behaviour that runs before and/or after it.
type Middleware func(http.Handler) http.Handler

// Chain wraps h so that mws run in the order given: the first is outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type userIDKey struct{}

// WithUserID returns a copy of ctx that carries the authenticated user's ID.
func WithUserID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

// UserID returns the authenticated user's ID, set by the JWT middleware.
func UserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey{}).(int)
	return id, ok
}

// QueryDefault returns the query parameter key, or def when it is absent.
// A parameter that is present but empty returns "".
func QueryDefault(r *http.Request, key, def string) string {
	q := r.URL.Query()
	if !q.Has(key) {
		return def
	}
	return q.Get(key)
}

// Route returns the matched route template without its method, for example
// "/api/reservations/{xid}", or "" when no route matched. ServeMux sets
// r.Pattern on the request it routes, so outer middleware can read it after
// calling the next handler.
func Route(r *http.Request) string {
	pattern := r.Pattern
	if i := strings.IndexByte(pattern, ' '); i >= 0 {
		pattern = pattern[i+1:]
	}
	return pattern
}

// ClientIP returns the caller's IP address. Proxy headers (X-Forwarded-For,
// then X-Real-IP) are only trusted when the direct peer is a loopback or
// private address, which is where the reverse proxy sits in every deployment
// of this API. A client on the public internet cannot spoof its address.
func ClientIP(r *http.Request) string {
	peer := remoteIP(r)
	if !trustedProxy(peer) {
		return peer
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		hops := strings.Split(xff, ",")
		// Walk from the nearest hop back; the first untrusted hop is the client.
		for i := len(hops) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(hops[i])
			if net.ParseIP(ip) == nil {
				break
			}
			if i == 0 || !trustedProxy(ip) {
				return ip
			}
		}
	}
	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(ip) != nil {
		return ip
	}
	return peer
}

func remoteIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func trustedProxy(ip string) bool {
	parsed := net.ParseIP(ip)
	return parsed != nil && (parsed.IsLoopback() || parsed.IsPrivate())
}

// StatusWriter records the status code and body size written through it.
type StatusWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// NewStatusWriter wraps w. If w is already a StatusWriter, it is returned as is.
func NewStatusWriter(w http.ResponseWriter) *StatusWriter {
	if sw, ok := w.(*StatusWriter); ok {
		return sw
	}
	return &StatusWriter{ResponseWriter: w}
}

func (w *StatusWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *StatusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

// Status is the response status code; 200 if the handler wrote nothing explicit.
func (w *StatusWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

// Written reports whether the handler has started the response.
func (w *StatusWriter) Written() bool { return w.status != 0 }

// Size is the number of body bytes written.
func (w *StatusWriter) Size() int { return w.size }

// Unwrap lets http.ResponseController reach the underlying writer.
func (w *StatusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
