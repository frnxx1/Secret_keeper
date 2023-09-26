package storage

import (
	"testing"
)

func TestDummyKeeper_Set(t *testing.T) {
	// смысла в тесте нет, так как ты проверяешь только заглушку
	keeper := GetDummyKeeper()
	key := "foo"
	value := "bar"

	keeper.Set(key, value, 0)
	if keeper.mem[key] != value {
		t.Error("bad memery value")
	}
}

func TestDummyKeeper_Get(t *testing.T) {
	// смысла в тесте нет, так как ты проверяешь только заглушку
	key := "foo"
	value := "bar"
	keeper := GetDummyKeeper()
	keeper.Set(key, value, 0)
	valueFromGet, _ := keeper.Get(key)
	if valueFromGet != value {
		t.Error("bad value from GET")
	}

}

func TestDummyKeeper_Clean(t *testing.T) {
	// смысла в тесте нет, так как ты проверяешь только заглушку
	key := "foo"
	keeper := GetDummyKeeper()
	valueFromGet := keeper.Clean(key)
	if valueFromGet != nil {
		t.Error("bad value from GET")
	}
}
