package storage

import (
	"sync"
	"testing"
)

func TestDummyKeeper_Set(t *testing.T) {
	keeper := DummyKeeper{Mem: make(map[string]string), mu: &sync.Mutex{}}
	key := "foo"
	value := "bar"
	
	keeper.Set(key, value,0)
	if keeper.Mem[key] != value {
		t.Error("bad memery value")
	}
}

func TestDummyKeeper_Get(t *testing.T) {
	key := "foo"
	value := "bar"
	keeper := DummyKeeper{Mem: make(map[string]string), mu: &sync.Mutex{}}
	keeper.Mem[key] = value
	value_from_get, _ := keeper.Get(key)
	if value_from_get != value {
		t.Error("bad value from GET")
	}

}

func TestDummyKeeper_Clean(t *testing.T) {
	key := "foo"
	keeper := DummyKeeper{Mem: make(map[string]string), mu: &sync.Mutex{}}
	value_from_get := keeper.Clean(key)
	if value_from_get != nil {
		t.Error("bad value from GET")
	}

}
