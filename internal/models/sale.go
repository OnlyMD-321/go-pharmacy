package models

import "time"

type Sale struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`      // Who made the sale
	InventoryID int64     `db:"inventory_id"` // Which item sold
	Quantity    int       `db:"quantity"`
	TotalPrice  float64   `db:"total_price"`
	SoldAt      time.Time `db:"sold_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
