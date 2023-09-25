package storage

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

const NotFoundError = "not found"

type Keeper interface {
	Get(key string) (string, error)
	Set(key, message string, ttl int) error
}

func GetDummyKeeper() Keeper {
	return DummyKeeper{Mem: make(map[string]string), mu: &sync.Mutex{}}
}

func GetRedisKeeper() Keeper {
	return RedisKeeper{*redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}), context.Background()}

}
