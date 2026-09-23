package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
	client *redis.Client
}

func NewRedisPublisher(addr string) *RedisPublisher {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisPublisher{
		client: rdb,
	}
}

func (r *RedisPublisher) PublishEvent(channel string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = r.client.Publish(context.Background(), channel, data).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to redis: %w", err)
	}

	return nil
}
