package cache

import (
	"testing"
	"time"
)

var cache *Cache

func init() {
	checkIntervalTime := 1 * time.Second
	holdInMemoryTime := 1 * time.Second
	cache = InitNewCache(checkIntervalTime, holdInMemoryTime)
}

func TestWrongInit1(t *testing.T) {
	defer func() {
		recover()
	}()
	InitNewCache(-1*time.Second, 1*time.Second)
	t.Error("Not panic on wrong parameters")
}

func TestWrongInit2(t *testing.T) {
	defer func() {
		recover()
	}()
	InitNewCache(1*time.Second, -1*time.Second)
	t.Error("Not panic on wrong parameters")
}

func TestGet(t *testing.T) {
	cache.Add("key", "test_data")
	data, _ := cache.Get("key")
	if data == nil {
		t.Errorf("Data not saved %s", data)
	}
}

func TestRewriting(t *testing.T) {
	cache.Add("key_to_rewrite", 40)
	cache.Add("key_to_rewrite", 42)
	data, _ := cache.Get("key_to_rewrite")
	if data != 42 {
		t.Errorf("Data don't rewriting")
	}
}

func TestDelete(t *testing.T) {
	cache.Add("key_to_delete", 42)
	err := cache.Delete("key_to_delete")
	if err != nil {
		t.Error("Error on delete ", err)
	}
	data, err := cache.Get("hey_to_delete")
	if err == nil {
		t.Error("Don't get error while getting deleted data")
	}
	if data != nil {
		t.Error("Error data", data)
	}
}

func TestCleaner(t *testing.T) {
	cache.Add("cleaner_key", 42)
	time.Sleep(2 * time.Second)
	_, err := cache.Get("cleaner key")
	if err == nil {
		t.Error("Cleaner don't delete data")
	}
}
