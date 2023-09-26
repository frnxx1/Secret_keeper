package storage

import (
	"errors"
	"sync"
)

type DummyKeeper struct {
	mem map[string]string // должно быть приватным
	mu  sync.Mutex
}

func GetDummyKeeper() *DummyKeeper {
	return &DummyKeeper{mem: make(map[string]string)}
}

func (k *DummyKeeper) Get(key string) (string, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	value, ok := k.mem[key]
	if !ok {
		return "", errors.New("not found")
	}

	k.Clean(key)

	return value, nil
}

func (k *DummyKeeper) Set(key, message string, _ int) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.mem[key] = message
	return nil
}

func (k *DummyKeeper) Clean(key string) error {
	// Clean публичный но нет под мютексом, это плохо.ы
	delete(k.mem, key)
	return nil
}
