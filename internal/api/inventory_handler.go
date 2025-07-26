package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OnlyMD-321/go-pharmacy/internal/models"
	"github.com/OnlyMD-321/go-pharmacy/internal/repositories"
)

type InventoryHandler struct {
	Repo *repositories.InventoryRepository
}

func NewInventoryHandler(db *pgxpool.Pool) *InventoryHandler {
	return &InventoryHandler{
		Repo: repositories.NewInventoryRepository(db),
	}
}

// GET /api/inventory
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	items, err := h.Repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /api/inventory
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req struct {
		Name        string    `json:"name" binding:"required"`
		Description string    `json:"description"`
		Quantity    int       `json:"quantity" binding:"required,min=0"`
		Price       float64   `json:"price" binding:"required,min=0"`
		ExpiryDate  time.Time `json:"expiry_date" time_format:"2006-01-02" time_utc:"1" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &models.InventoryItem{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		Price:       req.Price,
		ExpiryDate:  req.ExpiryDate,
	}

	err := h.Repo.Create(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}
