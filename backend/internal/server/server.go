package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"lvbp-project/backend/internal/auth"
	"lvbp-project/backend/internal/handlers/rest"
	"lvbp-project/backend/internal/handlers/sse"
	"lvbp-project/backend/internal/infrastructure/repositories"
	"lvbp-project/backend/internal/core/services"
	redisPkg "lvbp-project/backend/internal/pkg/redis"
)

type Server struct {
	router      *chi.Mux
	dbPool      *pgxpool.Pool
	redisClient *redisPkg.Client
}

func NewServer() *Server {
	// 1. Iniciar Redis
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisCfg := redisPkg.Config{
		Host:     redisHost,
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	redisClient, err := redisPkg.NewClient(redisCfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}
	fmt.Println("Successfully connected to Redis.")

	// 2. Iniciar PostgreSQL
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("Warning: Failed to connect to PostgreSQL: %v", err)
	} else {
		fmt.Println("Successfully connected to PostgreSQL.")
	}

	// 3. Inicializar Router base
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	srv := &Server{
		router:      r,
		dbPool:      dbPool,
		redisClient: redisClient,
	}

	// 4. Registrar módulos
	srv.registerModules()

	return srv
}

func (s *Server) registerModules() {
	auth.RegisterRoutes(s.router, s.dbPool, s.redisClient)

	// Dependency Injection for repositories and services
	gameRepo := repositories.NewGameRepository(s.dbPool)
	eventBus := redisPkg.NewEventBus(s.redisClient.GetDB())
	scorekeeperService := services.NewScorekeeperService(gameRepo, eventBus)

	// REST handlers
	boxscoreHandler := rest.NewBoxscoreHandler(gameRepo)
	ingestionHandler := rest.NewIngestionHandler(scorekeeperService)

	// SSE handler
	streamHandler := sse.NewStreamHandler(s.redisClient)

	// Mount /api routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// REST endpoints
		r.Get("/games/{gameId}/boxscore", boxscoreHandler.GetBoxscore)
		r.Post("/games/{gameId}/pitches", ingestionHandler.ProcessPitch)

		// SSE endpoint
		r.Get("/games/{gameId}/stream", streamHandler.StreamGameEvents)
	})
}


func (s *Server) Run(port string) error {
	defer func() {
		if s.dbPool != nil {
			s.dbPool.Close()
		}
		if s.redisClient != nil {
			s.redisClient.Close()
		}
	}()

	fmt.Printf("Server listening on port %s\n", port)
	return http.ListenAndServe(port, s.router)
}
