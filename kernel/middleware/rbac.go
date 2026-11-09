package middleware

import (
	"context"
	"labflow/httpx"
	"log"
	"net/http"
)

type RoleChecker interface {
	HasAccess(ctx context.Context, userID int, accessName string) (bool, error)
}

func RBACMiddleware(repo RoleChecker, permission string) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := httpx.UserID(r.Context())
			if !ok {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID not found"})
				return
			}
			allowed, err := repo.HasAccess(r.Context(), userID, permission)
			if err != nil {
				// Unknown is not "no": report it as a server problem, not a 403.
				log.Printf("[rbac] permission check failed for user %v (%s): %v", userID, permission, err)
				httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Permission check unavailable"})
				return
			}
			if !allowed {
				httpx.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "Permission denied"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
