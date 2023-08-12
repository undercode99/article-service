package caching

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/undercode99/article_service/config"
)

// NewRedisCaching creates a new instance of RedisCaching.
//
// It takes a pointer to a Config struct as its parameter and returns a pointer to a RedisCaching struct.
func NewRedisCaching(ctx context.Context, cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port),
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})

	// Check if the connection is successful
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalln("failed to connect to redis:", err)
	}

	return client
}
