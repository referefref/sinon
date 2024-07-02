package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func sendMetadataToRedis(config *Config, username, content, sourceIP string) error {
	if config.General.Redis.IP == "" || config.General.Redis.Port == 0 {
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.General.Redis.IP, config.General.Redis.Port),
	})
	ctx := context.Background()
	loginTime := time.Now().Format(time.RFC3339)
	data := fmt.Sprintf("Username: %s, Content: %s, SourceIP: %s, LoginTime: %s", username, content, sourceIP, loginTime)
	return rdb.Set(ctx, "session_metadata", data, 0).Err()
}
