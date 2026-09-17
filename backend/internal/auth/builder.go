package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	coreAuth "lvbp-project/backend/internal/auth/core"
	handlersAuth "lvbp-project/backend/internal/auth/handlers"
	pgInfra "lvbp-project/backend/internal/auth/infrastructure/postgres"
	redisInfra "lvbp-project/backend/internal/auth/infrastructure/redis"
	redisPkg "lvbp-project/backend/internal/pkg/redis"
)

// RegisterRoutes inicializa y registra todas las dependencias y rutas del módulo auth.
func RegisterRoutes(router chi.Router, dbPool *pgxpool.Pool, redisClient *redisPkg.Client) {
	// 1. Inicializar Repositorios
	userRepo := pgInfra.NewUserRepository(dbPool)
	sessionRepo := redisInfra.NewSessionRepository(redisClient)

	// 2. Inicializar Servicios
	jwtSecret := "super-secret-key-for-development-only"
	authService := coreAuth.NewAuthService(userRepo, sessionRepo, jwtSecret)

	// 3. Inicializar Handlers
	authHandler := handlersAuth.NewHandler(authService)

	// 4. Configurar Rutas
	router.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
	})
}
