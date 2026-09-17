package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository es un mock simple de UserRepository para las pruebas.
type mockUserRepo struct {
	user *User
	err  error
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return m.user, m.err
}

// MockSessionRepository es un mock simple de SessionRepository para las pruebas.
type mockSessionRepo struct {
	err error
}

func (m *mockSessionRepo) StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	return m.err
}

func (m *mockSessionRepo) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	return "", m.err
}

func (m *mockSessionRepo) RevokeToken(ctx context.Context, token string) error {
	return m.err
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()
	jwtSecret := "test-secret"
	validPassword := "mySuperSecretPassword123"

	// Generar el hash real utilizando bcrypt para simular el registro
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password for test setup: %v", err)
	}

	validUser := &User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockUser      *User
		mockUserErr   error
		mockStoreErr  error
		expectedError error
	}{
		{
			name:          "Success - Correct credentials",
			email:         "test@example.com",
			password:      validPassword,
			mockUser:      validUser,
			mockUserErr:   nil,
			mockStoreErr:  nil,
			expectedError: nil,
		},
		{
			name:          "Failure - Wrong password",
			email:         "test@example.com",
			password:      "wrongPassword",
			mockUser:      validUser,
			mockUserErr:   nil,
			mockStoreErr:  nil,
			expectedError: ErrInvalidCredentials,
		},
		{
			name:          "Failure - User not found",
			email:         "unknown@example.com",
			password:      "anyPassword",
			mockUser:      nil,
			mockUserErr:   errors.New("db record not found"),
			mockStoreErr:  nil,
			expectedError: ErrInvalidCredentials,
		},
		{
			name:          "Failure - Session store fails",
			email:         "test@example.com",
			password:      validPassword,
			mockUser:      validUser,
			mockUserErr:   nil,
			mockStoreErr:  errors.New("redis connection lost"),
			expectedError: errors.New("redis connection lost"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uRepo := &mockUserRepo{user: tt.mockUser, err: tt.mockUserErr}
			sRepo := &mockSessionRepo{err: tt.mockStoreErr}

			svc := NewAuthService(uRepo, sRepo, jwtSecret)

			accessToken, refreshToken, err := svc.Login(ctx, tt.email, tt.password)

			// Verificar si los errores coinciden
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
				// Asegurarnos de que no devolvió tokens en caso de error
				if accessToken != "" || refreshToken != "" {
					t.Errorf("Expected empty tokens on error, got Access: %v, Refresh: %v", accessToken, refreshToken)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if accessToken == "" {
					t.Errorf("Expected valid access token, got empty")
				}
				if refreshToken == "" {
					t.Errorf("Expected valid refresh token, got empty")
				}
			}
		})
	}
}
