package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Client representa la conexión a la base de datos Redis.
// Contiene la instancia subyacente del cliente go-redis.
type Client struct {
	rdb *redis.Client
}

// Config define los parámetros necesarios para conectarse a Redis.
type Config struct {
	Host     string // Dirección del servidor (ej. localhost)
	Port     string // Puerto de conexión (ej. 6379)
	Password string // Contraseña de autenticación (vacío si no tiene)
	DB       int    // Número de base de datos a utilizar
}

// NewClient inicializa y devuelve una nueva instancia de Client.
// Recibe un objeto Config y configura el cliente interno de Redis.
func NewClient(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Client{
		rdb: rdb,
	}, nil
}

// Ping verifica la conexión con el servidor Redis.
// Retorna un error si la base de datos no es alcanzable o la conexión falla.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.rdb.Ping(ctx).Result()
	return err
}

// Close cierra la conexión subyacente con Redis, liberando los recursos.
// Debe ser llamado (usualmente con defer) cuando el cliente ya no se necesite.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// GetDB expone el cliente interno de redis para comandos específicos de repositorios.
func (c *Client) GetDB() *redis.Client {
	return c.rdb
}
