package models

import "time"

type InventoryItem struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Quantity    int       `db:"quantity"`
	Price       float64   `db:"price"`
	ExpiryDate  time.Time `db:"expiry_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
