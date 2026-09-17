package redis

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create redis client: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client, got nil")
	}
	if client.rdb == nil {
		t.Fatal("Expected internal rdb to be initialized")
	}
}

// Note: This test requires a running Redis instance to pass.
// For unit tests, it's often mocked, but as requested we are doing small validation
// tests for inputs/outputs with the DB connection.
func TestPing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	cfg := Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create redis client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}
