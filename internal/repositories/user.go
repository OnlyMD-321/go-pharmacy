package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/OnlyMD-321/go-pharmacy/internal/models"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUID(ctx context.Context, uid string) (*models.User, error) {
	const query = `SELECT id, uid, name, email, role, created_at, updated_at FROM users WHERE uid=$1`

	user := &models.User{}
	err := r.DB.QueryRow(ctx, query, uid).Scan(
		&user.ID, &user.UID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
