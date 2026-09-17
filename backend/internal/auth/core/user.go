package auth

import "context"

// User representa el modelo de usuario en el dominio.
type User struct {
	ID           string
	Email        string
	PasswordHash string
}

// UserRepository define los métodos necesarios para interactuar con la persistencia de usuarios.
// Su implementación real (ej. PostgreSQL + sqlc) se inyectará al servicio.
type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
