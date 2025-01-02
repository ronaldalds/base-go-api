package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type InitRedis struct {
	Host     string
	Port     int
	Password string
	DB       int
}
type RedisStore struct {
	*redis.Client
}

func newRedisStore(init *InitRedis) *RedisStore {
	// Configura o cliente Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", init.Host, init.Port),
		Password: init.Password, // "" para Redis sem autenticação
		DB:       init.DB,       // Número do banco
	})

	// Verifica a conexão com o Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return &RedisStore{client}
}
func (rs *RedisStore) SetKey(key string, value string, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rs.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}
	return nil
}

func (rs *RedisStore) GetKey(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := rs.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", fmt.Errorf("failed to get key: %v", err)
	}
	return val, nil
}

func (rs *RedisStore) Health() map[string]string {
	stats := make(map[string]string)

	// Cria um contexto com timeout para o health check
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()

	// Testa a conectividade com o Redis usando PING
	_, err := rs.Client.Ping(ctx).Result()
	duration := time.Since(start)

	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to connect to Redis: %v", err)
	} else {
		stats["status"] = "up"
		stats["response_time"] = duration.String()
		stats["message"] = "Redis is healthy"
	}

	return stats
}
