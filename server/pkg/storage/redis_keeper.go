package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisKeeper struct {
	cn *redis.Client // вот тут лучше ссылку на клиент хранить, чем разименовывать клиент
	// ctx context.Context // хранить контекст плохая практика, его лучше всегда пробрасывать
}

// const TTL = 0 // не используется

func (k *RedisKeeper) Get(ctx context.Context, key string) (string, error) {
	val, err := k.cn.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(NotFoundError)
	}

	return val, err
}

func (k *RedisKeeper) Set(ctx context.Context, key, message string, ttl int) error {
	seconds, err := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	if err != nil {
		return err
	}

	return k.cn.Set(ctx, key, message, seconds).Err()
}
