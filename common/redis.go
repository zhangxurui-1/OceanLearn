package common

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var client *redis.Client

func InitRedis() *redis.Client {
	// 创建 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host" + ":" + viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	// 检查是否连接成功
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("failed to connect redis: %v", err)
		return nil
	}

	return redisClient
}

func GetRedisClient() *redis.Client {
	if client == nil {
		client = InitRedis()
	}
	return client
}
