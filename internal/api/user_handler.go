package api

import (
	"net/http"

	"github.com/OnlyMD-321/go-pharmacy/internal/middlewares"
	"github.com/OnlyMD-321/go-pharmacy/internal/models"
	"github.com/OnlyMD-321/go-pharmacy/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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

// POST /api/register
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		UID   string `json:"uid" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
		Role  string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Firebase UID from context (validated by middleware)
	firebaseUID, exists := c.Get(middlewares.ContextFirebaseUID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No Firebase UID found in context"})
		return
	}
	if firebaseUID != req.UID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UID mismatch between token and request"})
		return
	}

	// RBAC: Only allow certain roles to register users (e.g., only admin can create admin)
	// Example: Only allow self-registration as 'seller', others require admin
	if req.Role != "seller" {
		// Get current user from DB
		currentUser, _ := h.Repo.FindByUID(c.Request.Context(), req.UID)
		if currentUser == nil || currentUser.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create non-seller users"})
			return
		}
	}

	// Check if user exists
	existing, _ := h.Repo.FindByUID(c.Request.Context(), req.UID)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user := &models.User{
		UID:   req.UID,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}
	err := h.Repo.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}
