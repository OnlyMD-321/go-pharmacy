package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OnlyMD-321/go-pharmacy/internal/models"
)

type SaleRepository struct {
	DB *pgxpool.Pool
}

func NewSaleRepository(db *pgxpool.Pool) *SaleRepository {
	return &SaleRepository{DB: db}
}

func (r *SaleRepository) Create(ctx context.Context, sale *models.Sale) error {
	const query = `
	INSERT INTO sales (user_id, inventory_id, quantity, total_price, sold_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	now := time.Now().UTC()
	err := r.DB.QueryRow(ctx, query,
		sale.UserID,
		sale.InventoryID,
		sale.Quantity,
		sale.TotalPrice,
		sale.SoldAt,
		now,
		now,
	).Scan(&sale.ID)

	return err
}

func (r *SaleRepository) GetAll(ctx context.Context) ([]models.Sale, error) {
	const query = `SELECT id, user_id, inventory_id, quantity, total_price, sold_at, created_at, updated_at FROM sales`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var s models.Sale
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.InventoryID,
			&s.Quantity,
			&s.TotalPrice,
			&s.SoldAt,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}
	return sales, nil
}
