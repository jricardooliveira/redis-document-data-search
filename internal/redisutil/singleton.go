package redisutil

import (
	"sync"
	"github.com/redis/go-redis/v9"
)

var (
	singletonClient *redis.Client
	singletonOnce sync.Once
)

// GetSingletonRedisClient returns a singleton Redis client for the given URL (first call wins).
func GetSingletonRedisClient(redisURL string) *redis.Client {
	singletonOnce.Do(func() {
		client, err := NewRedisClient(redisURL)
		if err != nil {
			panic("Failed to create singleton Redis client: " + err.Error())
		}
		singletonClient = client
	})
	return singletonClient
}
