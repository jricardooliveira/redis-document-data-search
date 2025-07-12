package redisutil

import (
	"context"
	"fmt"
	"strings"
	"github.com/redis/go-redis/v9"
)

// GetRedisMemoryInfo returns used_memory and used_memory_human from INFO MEMORY
func GetRedisMemoryInfo(client *redis.Client) (usedBytes int64, usedHuman string, err error) {
	ctx := context.Background()
	info, err := client.Info(ctx, "memory").Result()
	if err != nil {
		return 0, "", err
	}
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "used_memory:") {
			_, val, _ := strings.Cut(line, ":")
			fmtVal := strings.TrimSpace(val)
			// parse int64
			var mem int64
			_, _ = fmt.Sscanf(fmtVal, "%d", &mem)
			usedBytes = mem
		}
		if strings.HasPrefix(line, "used_memory_human:") {
			_, val, _ := strings.Cut(line, ":")
			usedHuman = strings.TrimSpace(val)
		}
	}
	return usedBytes, usedHuman, nil
}
