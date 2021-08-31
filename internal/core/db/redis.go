package db

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedis(address, password string, db int) (redClient *RedisClient, err error) {

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	ctx := context.TODO()
	_, err = client.Ping(ctx).Result()

	redClient = &RedisClient{Client: client}

	return
}
