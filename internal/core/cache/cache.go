package cache

import (
	"context"
	"encoding/json"
	"gopher/internal/core/db"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	Cache *redis.Client
}

func RedisCache(redisCli *db.RedisClient) (redisCache *Client, err error) {
	redisCache = &Client{Cache: redisCli.Client}
	return
}

// set value by key
func (r *Client) Set(key string, value interface{}, ttl int) error {
	parseValue, err := json.Marshal(value)
	if err != nil {
		parseValue = []byte(value.(string))
	}
	ctx := context.TODO()
	return r.Cache.Set(ctx, key, string(parseValue), time.Duration(ttl)*time.Second).Err()
}

// get value by key
func (r *Client) Get(key string) (string, error) {
	ctx := context.TODO()
	value, err := r.Cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

// Increase use to counter
func (r *Client) Increase(key string) (int64, error) {
	ctx := context.TODO()
	val, err := r.Cache.Incr(ctx, key).Result()
	return val, err
}

// check to this key exist or not
func (r *Client) KeyExist(key string) (exist bool, err error) {
	var existInt int64
	ctx := context.TODO()
	existInt, err = r.Cache.Exists(ctx, key).Result()
	exist = (existInt == 1)
	return exist, err
}

// delete value by key
func (r *Client) Delete(key string) (err error) {
	ctx := context.TODO()
	err = r.Cache.Del(ctx, key).Err()
	if err != nil {
		return
	}
	return
}

// delete and update set of items by one key, like array of resource
func (r *Client) UpdateSet(key string, values []string, ttl int) (err error) {

	ctx := context.TODO()
	err = r.Cache.Del(ctx, key).Err()
	if err != nil {
		return
	}

	for _, v := range values {
		err = r.Cache.SAdd(ctx, key, v).Err()
		if err != nil {
			return
		}
	}

	err = r.Cache.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	return
}

// check to value member of set
func (r *Client) IsMemberSet(key, value string) (exist bool, err error) {
	ctx := context.TODO()
	exist, err = r.Cache.SIsMember(ctx, key, value).Result()
	return
}
