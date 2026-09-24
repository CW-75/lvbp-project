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

	// Dependency Injection for handlers
	// Currently GameRepository and StandingsRepository might need real implementations, 
	// here we just initialize the handlers to bind the routes if they were created.
	// We will create the routes inside the chi router.

	// For the sake of the project architecture, we'll mount /api
	s.router.Route("/api/v1", func(r chi.Router) {
		// // TODO: inject actual repositories
		// gameRepo := postgres.NewGameRepository(s.dbPool)
		// standingsRepo := postgres.NewStandingsRepository(s.dbPool)
		
		// restHandlers
		// boxscoreHandler := rest.NewBoxscoreHandler(gameRepo)
		// standingsHandler := rest.NewStandingsHandler(standingsRepo)
		
		// r.Get("/games/{gameId}/boxscore", boxscoreHandler.GetBoxscore)
		// r.Get("/standings", standingsHandler.GetStandings)
		
		// sse handler
		// streamHandler := sse.NewStreamHandler(s.redisClient)
		// r.Get("/games/{gameId}/stream", streamHandler.StreamGameEvents)
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
