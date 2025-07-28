package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OnlyMD-321/go-pharmacy/internal/repositories"
)

type RBACMiddleware struct {
	UserRepo *repositories.UserRepository
	AllowedRoles []string
}

// NewRBACMiddleware returns a Gin middleware handler that authorizes based on roles.
func NewRBACMiddleware(db *pgxpool.Pool, allowedRoles ...string) gin.HandlerFunc {
	userRepo := repositories.NewUserRepository(db)

	return func(c *gin.Context) {
		uidValue, exists := c.Get(ContextFirebaseUID)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Firebase UID not found in context"})
			return
		}
		uid, ok := uidValue.(string)
		if !ok || strings.TrimSpace(uid) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Firebase UID"})
			return
		}

		// Fetch user from DB
		user, err := userRepo.FindByUID(context.Background(), uid)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User not found or unauthorized"})
			return
		}

		// Check if user role is allowed
		for _, role := range allowedRoles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient permissions"})
	}
}
