package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OnlyMD-321/go-pharmacy/internal/middlewares"
	"github.com/OnlyMD-321/go-pharmacy/internal/repositories"
)

type UserHandler struct {
	Repo *repositories.UserRepository
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		Repo: repositories.NewUserRepository(db),
	}
}

// GET /api/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	uidValue, exists := c.Get(middlewares.ContextFirebaseUID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No Firebase UID found in context"})
		return
	}
	uid, ok := uidValue.(string)
	if !ok || uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Firebase UID"})
		return
	}

	user, err := h.Repo.FindByUID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
