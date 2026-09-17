package redis

import (
	"context"
	"fmt"
	"time"
	
	"lvbp-project/backend/internal/pkg/redis"
)

// SessionRepository maneja el almacenamiento de los tokens de sesión (refresh tokens) en Redis.
type SessionRepository struct {
	client *redis.Client
}

// NewSessionRepository crea una nueva instancia del repositorio de sesiones.
func NewSessionRepository(client *redis.Client) *SessionRepository {
	return &SessionRepository{
		client: client,
	}
}

// StoreRefreshToken almacena un refresh token en Redis con un TTL (Time To Live).
// El token expirará automáticamente una vez cumplido el TTL.
func (r *SessionRepository) StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	// Guardamos el userID como valor asociado al token para poder identificar a quién pertenece al renovar.
	err := r.client.GetDB().Set(ctx, key, userID, ttl).Err()
	return err
}

// GetUserIDByToken obtiene el ID del usuario asociado a un refresh token.
// Retorna error si el token no existe (expiró) o es inválido.
func (r *SessionRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("session:%s", token)
	userID, err := r.client.GetDB().Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return userID, nil
}

// RevokeToken elimina un refresh token de Redis, cerrando la sesión efectivamente.
func (r *SessionRepository) RevokeToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	err := r.client.GetDB().Del(ctx, key).Err()
	return err
}
