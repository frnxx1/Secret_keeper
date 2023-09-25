package storage

import (
	"errors"
	"sync"
	
)

type DummyKeeper struct {
	Mem map[string]string
	mu  *sync.Mutex
}

func (k DummyKeeper) Get(key string) (string, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	value, ok := k.Mem[key]
	if !ok {

		return "", errors.New("not found")
	}
	
	k.Clean(key)

	return value, nil
}

func (k DummyKeeper) Set(key, message string,ttl int) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.Mem[key] = message
	return nil
}

func (k DummyKeeper) Clean(key string) error {

	delete(k.Mem, key)
	return nil
}
