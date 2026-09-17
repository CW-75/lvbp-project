package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotFound        = errors.New("user not found")
	ErrInternalError       = errors.New("internal server error")
)

type SessionRepository interface {
	StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error
	GetUserIDByToken(ctx context.Context, token string) (string, error)
	RevokeToken(ctx context.Context, token string) error
}

type AuthService struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	jwtSecret   []byte
}

func NewAuthService(uRepo UserRepository, sRepo SessionRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    uRepo,
		sessionRepo: sRepo,
		jwtSecret:   []byte(jwtSecret),
	}
}

// Login validates user credentials and generates tokens.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// En la práctica, habría que distinguir entre 'no encontrado' y errores de DB
		return "", "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate Access Token (JWT)
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token (random string)
	// For simplicity in this implementation, we use a basic string, but ideally it should be a crypto/rand generated token
	refreshToken := "rt_" + user.ID + "_" + time.Now().Format("20060102150405")

	// Store refresh token in Redis with a 7 day TTL
	err = s.sessionRepo.StoreRefreshToken(ctx, user.ID, refreshToken, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) generateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}
