package store

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/whj1990/go-core/config"
	"github.com/whj1990/go-core/trace"
	"time"
)

type RedisClient struct {
	*redis.Client
}

func (r *RedisClient) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	marshalObject, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.Set(ctx, key, marshalObject, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetObject(ctx context.Context, key string, obj interface{}) error {
	val, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(val, obj); err != nil {
		return err
	}
	return nil
}

func NewRedisClient() *RedisClient {
	client := &RedisClient{
		redis.NewClient(&redis.Options{
			Addr:     config.GetNaCosString("redis.address", ""),
			Password: config.GetNaCosString("redis.password", ""),
			DB:       config.GetNaCosInt("redis.database", 0),
		}),
	}
	client.AddHook(trace.NewTracingHook())
	return client
}
