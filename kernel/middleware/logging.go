package middleware

import (
	"fmt"
	"io"
	"labflow/httpx"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// AccessLog writes one line per request to out, in the same shape as the
// Gin access log this API used before: time | status | latency | ip | method path.
func AccessLog(out io.Writer) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := httpx.NewStatusWriter(w)
			next.ServeHTTP(sw, r)
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path += "?" + r.URL.RawQuery
			}
			fmt.Fprintf(out, "[http] %s | %3d | %13v | %15s | %-7s %q\n",
				start.Format("2006/01/02 - 15:04:05"), sw.Status(), time.Since(start),
				httpx.ClientIP(r), r.Method, path)
		})
	}
}

// Recover turns a panic in a handler into a 500 response and logs the stack.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := httpx.NewStatusWriter(w)
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if rec == http.ErrAbortHandler {
				panic(rec)
			}
			log.Printf("[recover] panic serving %s %s: %v\n%s", r.Method, r.URL.Path, rec, debug.Stack())
			if !sw.Written() {
				sw.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(sw, r)
	})
}
