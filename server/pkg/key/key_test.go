package key

import "testing"

func TestDummyKeyBuilder(t *testing.T) {
	// тоже нет смысла в тесте
	dummyKeyBuilder := DummyKeyBuilder{}
	key, _ := dummyKeyBuilder.Get()
	if key != DummyTestKey {
		t.Error("bad dummy key")
	}
}
