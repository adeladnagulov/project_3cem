package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCnf struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisDb(cnf RedisCnf) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cnf.Addr,
		Password:     cnf.Password,
		DB:           cnf.DB,
		PoolSize:     20,
		MinIdleConns: 10,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis %w", err)
	}
	return client, nil
}
