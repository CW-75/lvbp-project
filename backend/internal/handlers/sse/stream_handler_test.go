package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	redisPkg "lvbp-project/backend/internal/pkg/redis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupRedisContainer levanta un contenedor de Redis para la prueba E2E del SSE.
func setupRedisContainer(ctx context.Context, t *testing.T) (*redisPkg.Client, func()) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := redisC.Host(ctx)
	require.NoError(t, err)
	port, err := redisC.MappedPort(ctx, "6379")
	require.NoError(t, err)

	client, err := redisPkg.NewClient(redisPkg.Config{
		Host: host,
		Port: port.Port(),
	})
	require.NoError(t, err)

	return client, func() {
		client.Close()
		redisC.Terminate(ctx)
	}
}

func TestStreamGameEvents(t *testing.T) {
	ctx := context.Background()
	client, teardown := setupRedisContainer(ctx, t)
	defer teardown()

	handler := NewStreamHandler(client)

	// Configurar Request con Chi context
	req, err := http.NewRequest("GET", "/games/test-game/stream", nil)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("gameId", "test-game")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Como httptest.ResponseRecorder expone http.Flusher de forma implícita, podemos usarlo
	rr := httptest.NewRecorder()

	// Lanzar la respuesta del stream en una goroutine
	done := make(chan struct{})
	go func() {
		handler.StreamGameEvents(rr, req)
		close(done)
	}()

	// Esperamos a que inicie la conexión y publique el evento 'connected'
	time.Sleep(100 * time.Millisecond)

	// Publicamos un evento simulado al canal de Redis
	err = client.GetDB().Publish(ctx, "game:test-game:events", `{"type":"PITCH","result":"STRIKE"}`).Err()
	assert.NoError(t, err)

	// Esperamos para asegurar que SSE lo consuma
	time.Sleep(100 * time.Millisecond)

	// Verificamos que los headers se asignaron correctamente
	assert.Equal(t, "text/event-stream", rr.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", rr.Header().Get("Cache-Control"))

	body := rr.Body.String()
	assert.Contains(t, body, "event: connected")
	assert.Contains(t, body, `{"type":"PITCH","result":"STRIKE"}`)
}
