package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"hezzlService/src/pkg/config"
	"log"
	"time"
)

type RedisClient struct {
	key           string
	writeDuration time.Duration
	*redis.Client
}

func NewRedisClient(cfg config.Cache) *RedisClient {
	return &RedisClient{
		cfg.Secret,
		time.Duration(cfg.WriteDuration) * time.Second,
		redis.NewClient(&redis.Options{
			Addr:     cfg.Host + cfg.Port,
			Password: cfg.Pwd,
			DB:       0,
		}),
	}
}

func (r *RedisClient) GetItems(ctx context.Context) []byte {
	data, err := r.Get(ctx, r.key).Result()
	if err != nil {
		log.Printf("[REDIS] | GET items/list | failed: %s\n", err.Error())
		return nil
	}

	return []byte(data)
}

func (r *RedisClient) Save(ctx context.Context, items []byte) error {
	err := r.Set(ctx, r.key, items, r.writeDuration).Err()
	if err != nil {
		return err
	}

	return nil
}
