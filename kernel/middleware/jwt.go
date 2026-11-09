package middleware

import (
	"labflow/auth"
	"labflow/httpx"
	"net/http"
	"strings"
)

func JWTMiddleware(tokenUtil *auth.TokenUtil) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized - No token provided"})
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			// 2. Extract token from bearer
			if tokenString == authHeader {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized - Invalid token format"})
				return
			}

			// 3. Pass and validate the token
			authToken, err := tokenUtil.ParseJWT(r.Context(), tokenString)
			if err != nil {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized - Invalid token"})
				return
			}
			next.ServeHTTP(w, r.WithContext(httpx.WithUserID(r.Context(), authToken.ID)))
		})
	}
}

func JWTCookieMiddleware(tokenUtil *auth.TokenUtil) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized - No token provided"})
				return
			}

			// 3. Pass and validate the token
			authToken, err := tokenUtil.ParseJWT(r.Context(), cookie.Value)
			if err != nil {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized - Invalid token"})
				return
			}
			next.ServeHTTP(w, r.WithContext(httpx.WithUserID(r.Context(), authToken.ID)))
		})
	}
}
