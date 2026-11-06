package middleware

import (
	"labflow/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(tokenUtil *auth.TokenUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token provided"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// 2. Extract token from bearer
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token format"})
			c.Abort()
			return
		}

		// 3. Pass and validate the token
		authToken, err := tokenUtil.ParseJWT(c, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token"})
			c.Abort()
			return
		}
		c.Set("id", authToken.ID)
		c.Next()
	}
}

func JWTCookieMiddleware(tokenUtil *auth.TokenUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token provided"})
			c.Abort()
			return
		}

		// 3. Pass and validate the token
		authToken, err := tokenUtil.ParseJWT(c, authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token"})
			c.Abort()
			return
		}
		c.Set("id", authToken.ID)
		c.Next()
	}
}
