package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"lvbp-project/backend/internal/auth/core"
)

// UserRepository implementa auth.UserRepository utilizando PostgreSQL.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository inicializa un nuevo repositorio de usuarios.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserByEmail busca un usuario por su correo electrónico.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	// TODO: Esta consulta será reemplazada por sqlc en el futuro.
	query := `SELECT id, email, password_hash FROM users WHERE email = $1 LIMIT 1`
	
	var user auth.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}
