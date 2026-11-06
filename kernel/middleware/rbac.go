package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleChecker interface {
	HasAccess(ctx context.Context, userID int, accessName string) (bool, error)
}

func RBACMiddleware(repo RoleChecker, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			c.Abort()
			return
		}
		allowed, err := repo.HasAccess(c.Request.Context(), userID.(int), permission)
		if err != nil {
			// Unknown is not "no": report it as a server problem, not a 403.
			log.Printf("[rbac] permission check failed for user %v (%s): %v", userID, permission, err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Permission check unavailable"})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
