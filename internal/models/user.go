package models

import "time"

type User struct {
	ID        int64     `db:"id"`
	UID       string    `db:"uid"`    // Firebase UID
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`   // admin, pharmacist, seller
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
