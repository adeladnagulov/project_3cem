package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("not found")

type RedisCodesRepo struct {
	RC     *redis.Client
	prefix string
	ttl    time.Duration
}

func NewRedesRepoCodes(client *redis.Client) *RedisCodesRepo {
	return &RedisCodesRepo{
		RC:     client,
		prefix: "confirm:",
		ttl:    1 * time.Minute,
	}
}

func (r *RedisCodesRepo) AddNewCode(code string, creator string) error {
	return r.RC.Set(context.Background(), r.newKey(creator), code, r.ttl).Err()
}

func (r *RedisCodesRepo) ValidateCode(code string, email string) (bool, error) {
	value, err := r.RC.Get(context.Background(), r.newKey(email)).Result()
	if err == redis.Nil {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	return value == code, nil
}

func (r *RedisCodesRepo) newKey(creator string) string {
	return r.prefix + creator
}
