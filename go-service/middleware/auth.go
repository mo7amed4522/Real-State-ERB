package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-Id")
		userType := c.GetHeader("X-User-Type")
		if userIDStr == "" || userType == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication headers"})
			return
		}
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
			return
		}
		ctx := c.Request.Context()
		ctx = contextWithUser(ctx, uint(userID), userType)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Helper to inject user info into context
import "context"
type contextKey string

func contextWithUser(ctx context.Context, userID uint, userType string) context.Context {
	ctx = context.WithValue(ctx, contextKey("user_id"), userID)
	ctx = context.WithValue(ctx, contextKey("user_type"), userType)
	return ctx
} 