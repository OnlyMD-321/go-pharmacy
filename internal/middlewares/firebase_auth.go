package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/OnlyMD-321/go-pharmacy/internal/firebase"
	"firebase.google.com/go/v4/auth"
)

const ContextFirebaseUID = "firebaseUID"

// FirebaseAuthMiddleware validates Firebase ID token and adds UID to context
func FirebaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// Expecting header: "Bearer <token>"
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader { // no "Bearer " prefix found
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		client, err := firebase.App().Auth(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firebase Auth client"})
			return
		}

		token, err := client.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired Firebase ID token"})
			return
		}

		// Add Firebase UID to Gin context
		c.Set(ContextFirebaseUID, token.UID)

		c.Next()
	}
}
