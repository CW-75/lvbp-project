package pubsub_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"

	"lvbp-project/backend/internal/infrastructure/pubsub"
)

func TestRedisPublisher_PublishEvent(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	endpoint, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Clean up "redis://" prefix since go-redis expects host:port
	addr := endpoint
	if len(addr) > 8 && addr[:8] == "redis://" {
		addr = addr[8:]
	}

	publisher := pubsub.NewRedisPublisher(addr)

	type testEvent struct {
		Action string
		Data   int
	}

	evt := testEvent{
		Action: "PitchThrown",
		Data:   100,
	}

	err = publisher.PublishEvent("game:123", evt)
	assert.NoError(t, err)
	
	// A proper integration test would also create a subscriber and verify the message is received,
	// but this ensures the publisher can connect and marshal the data without error.
	
	// Add small delay to ensure completion
	time.Sleep(100 * time.Millisecond)
}
