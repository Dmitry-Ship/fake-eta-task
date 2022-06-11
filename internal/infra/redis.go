package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func getRedisClient(ctx context.Context) *redis.Client {
	port := os.Getenv("REDIS_PORT")
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	options := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     options,
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	log.Printf("📢 Connected to redis %s", options)

	return client
}

type Cache interface {
	Get(key string, result interface{}) error
	Set(key string, value interface{}) error
}

type cache struct {
	client *redis.Client
	ctx    context.Context
}

func NewCache(ctx context.Context) *cache {
	return &cache{client: getRedisClient(ctx), ctx: ctx}
}

func (c *cache) Get(key string, result interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), result)
}

func (c *cache) Set(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, val, time.Second*10).Err()
}
