package storage

import (
	"github.com/redis/go-redis/v9"
)

const NotFoundError = "not found"

func GetRedisKeeper() *RedisKeeper {
	return &RedisKeeper{
		cn: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // такие вещи лучше в коде не хардкодить а пробрасывать через конфиг, параметры, переменные окружения итп...
			Password: "",
			DB:       0,
		}),
	}
}
