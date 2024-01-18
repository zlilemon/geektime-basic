package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"geektime-basic/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	Key(id int64) string
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.Key(id)

	// 数据不存在, err = redis.Nil
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	err = json.Unmarshal(val, &u)

	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}

	key := cache.Key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) Key(id int64) string {
	//return fmt.Sprintf("user:info:#{id}")
	return fmt.Sprintf("user:info:%d", id)
}
