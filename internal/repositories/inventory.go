package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OnlyMD-321/go-pharmacy/internal/models"
)

type InventoryRepository struct {
	DB *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{DB: db}
}

func (r *InventoryRepository) Create(ctx context.Context, item *models.InventoryItem) error {
	const query = `
	INSERT INTO inventory (name, description, quantity, price, expiry_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	now := time.Now().UTC()
	err := r.DB.QueryRow(ctx, query,
		item.Name,
		item.Description,
		item.Quantity,
		item.Price,
		item.ExpiryDate,
		now,
		now,
	).Scan(&item.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *InventoryRepository) GetAll(ctx context.Context) ([]models.InventoryItem, error) {
	const query = `SELECT id, name, description, quantity, price, expiry_date, created_at, updated_at FROM inventory`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.InventoryItem{}
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Quantity,
			&item.Price,
			&item.ExpiryDate,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
