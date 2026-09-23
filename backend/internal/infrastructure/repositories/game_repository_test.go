package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"lvbp-project/backend/internal/core/domain"
	"lvbp-project/backend/internal/infrastructure/repositories"
)

func setupTestDB(ctx context.Context, t *testing.T) (*pgxpool.Pool, func()) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts("../../../sql/schema.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	dbPool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	cleanup := func() {
		dbPool.Close()
		pgContainer.Terminate(ctx)
	}

	return dbPool, cleanup
}

func TestGameRepository_CreateAndGetGame(t *testing.T) {
	ctx := context.Background()
	dbPool, cleanup := setupTestDB(ctx, t)
	defer cleanup()

	// Seed teams first
	_, err := dbPool.Exec(ctx, `INSERT INTO teams (id, name, short_name) VALUES ('CAR', 'Leones del Caracas', 'Caracas'), ('MAG', 'Navegantes del Magallanes', 'Magallanes')`)
	require.NoError(t, err)

	repo := repositories.NewGameRepository(dbPool)

	// Create Game
	newGame := &domain.Game{
		HomeTeamID: "CAR",
		AwayTeamID: "MAG",
		StartTime:  time.Now(),
	}

	createdGame, err := repo.CreateGame(ctx, newGame)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, createdGame.ID)
	assert.Equal(t, "CAR", createdGame.HomeTeamID)

	// Get Game
	retrievedGame, err := repo.GetGameByID(ctx, createdGame.ID)
	require.NoError(t, err)
	assert.Equal(t, createdGame.ID, retrievedGame.ID)
	assert.Equal(t, domain.GameStatusScheduled, retrievedGame.Status)
}

func TestGameRepository_UpdateGame(t *testing.T) {
	ctx := context.Background()
	dbPool, cleanup := setupTestDB(ctx, t)
	defer cleanup()

	_, err := dbPool.Exec(ctx, `INSERT INTO teams (id, name, short_name) VALUES ('CAR', 'Leones', 'CAR'), ('MAG', 'Navegantes', 'MAG')`)
	require.NoError(t, err)

	repo := repositories.NewGameRepository(dbPool)

	game, err := repo.CreateGame(ctx, &domain.Game{
		HomeTeamID: "CAR",
		AwayTeamID: "MAG",
		StartTime:  time.Now(),
	})
	require.NoError(t, err)

	// Update game
	game.Status = domain.GameStatusInProgress
	game.HomeScore = 2
	game.CurrentInning = 3
	game.IsTopInning = false

	updated, err := repo.UpdateGame(ctx, game)
	require.NoError(t, err)
	assert.Equal(t, domain.GameStatusInProgress, updated.Status)
	assert.Equal(t, 2, updated.HomeScore)
	assert.Equal(t, 3, updated.CurrentInning)
	assert.False(t, updated.IsTopInning)
}
