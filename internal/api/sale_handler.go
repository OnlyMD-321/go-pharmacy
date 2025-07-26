package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OnlyMD-321/go-pharmacy/internal/models"
	"github.com/OnlyMD-321/go-pharmacy/internal/repositories"
)

type SaleHandler struct {
	Repo *repositories.SaleRepository
}

func NewSaleHandler(db *pgxpool.Pool) *SaleHandler {
	return &SaleHandler{
		Repo: repositories.NewSaleRepository(db),
	}
}

// GET /api/sales
func (h *SaleHandler) GetSales(c *gin.Context) {
	sales, err := h.Repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales"})
		return
	}
	c.JSON(http.StatusOK, sales)
}

// POST /api/sales
func (h *SaleHandler) CreateSale(c *gin.Context) {
	var req struct {
		UserID      int64     `json:"user_id" binding:"required"`
		InventoryID int64     `json:"inventory_id" binding:"required"`
		Quantity    int       `json:"quantity" binding:"required,min=1"`
		TotalPrice  float64   `json:"total_price" binding:"required,min=0"`
		SoldAt      time.Time `json:"sold_at" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sale := &models.Sale{
		UserID:      req.UserID,
		InventoryID: req.InventoryID,
		Quantity:    req.Quantity,
		TotalPrice:  req.TotalPrice,
		SoldAt:      req.SoldAt,
	}

	if err := h.Repo.Create(c.Request.Context(), sale); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sale"})
		return
	}

	c.JSON(http.StatusCreated, sale)
}
