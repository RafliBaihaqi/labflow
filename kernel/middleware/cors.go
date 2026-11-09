package middleware

import (
	"labflow/httpx"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// CORS answers preflight requests and adds CORS headers for allowed origins.
// A cross-origin request from an origin that is not allowed gets 403.
// Requests without an Origin header, or from the API's own host, pass through.
func CORS(cfg CORSConfig) httpx.Middleware {
	normal := http.Header{}
	preflight := http.Header{}
	if cfg.AllowCredentials {
		normal.Set("Access-Control-Allow-Credentials", "true")
		preflight.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(cfg.ExposeHeaders) > 0 {
		normal.Set("Access-Control-Expose-Headers", joinHeaders(cfg.ExposeHeaders, http.CanonicalHeaderKey))
	}
	normal.Set("Vary", "Origin")
	if len(cfg.AllowMethods) > 0 {
		preflight.Set("Access-Control-Allow-Methods", joinHeaders(cfg.AllowMethods, strings.ToUpper))
	}
	if len(cfg.AllowHeaders) > 0 {
		preflight.Set("Access-Control-Allow-Headers", joinHeaders(cfg.AllowHeaders, http.CanonicalHeaderKey))
	}
	if cfg.MaxAge > 0 {
		preflight.Set("Access-Control-Max-Age", strconv.FormatInt(int64(cfg.MaxAge/time.Second), 10))
	}
	preflight.Add("Vary", "Origin")
	preflight.Add("Vary", "Access-Control-Request-Method")
	preflight.Add("Vary", "Access-Control-Request-Headers")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" || origin == "http://"+r.Host || origin == "https://"+r.Host {
				next.ServeHTTP(w, r)
				return
			}
			if !slices.Contains(cfg.AllowOrigins, origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			headers := normal
			if r.Method == http.MethodOptions {
				headers = preflight
			}
			for key, values := range headers {
				w.Header()[key] = values
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func joinHeaders(values []string, normalize func(string) string) string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, normalize(v))
		}
	}
	return strings.Join(out, ",")
}
