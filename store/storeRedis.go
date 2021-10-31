package store

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/xsoroton/autopilot-proxy/util"
)

// REDIS implementation of store
type RedisStore struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedisStoreFromEnv() Store {
	client := redis.NewClient(&redis.Options{
		Addr:     util.GetEnv("REDIS_HOST", "localhost:6379"),
		Password: "",
		DB:       0,
	})
	return &RedisStore{
		Client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisStore) Set(key string, value []byte) (err error) {
	err = r.Client.Set(r.ctx, key, value, 0).Err()
	return
}

func (r *RedisStore) Get(key string) (value []byte, err error) {
	value, err = r.Client.Get(r.ctx, key).Bytes()
	return
}

func (r *RedisStore) Remove(key string) (err error) {
	err = r.Client.Del(r.ctx, key).Err()
	return
}
