package redis_client

import (
	"context"
	"fmt"
	"livebets/analazer/cmd/config"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		// Password: cfg.Password,
		DB: int(cfg.DB)})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) Publish(ctx context.Context, key string, value interface{}) error {
	err := r.client.Publish(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) PSubscribe(ctx context.Context, pattern string) *redis.PubSub {
	return r.client.PSubscribe(ctx, pattern)
}

func (r Redis) Close() error {
	return r.client.Close()
}
