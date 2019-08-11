package cache

import (
	"testing"
	"time"
)

var cache *Cache
var fastCache *Cache

func init() {
	checkIntervalTime := 1 * time.Second
	holdInMemoryTime := 1 * time.Second
	cache = InitNewCache(checkIntervalTime, holdInMemoryTime)

	shortCheckIntervalTime := 1 * time.Millisecond
	shortHoldInMemoryTime := 1 * time.Millisecond
	fastCache = InitNewCache(shortCheckIntervalTime, shortHoldInMemoryTime)
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
	key := "key"
	cache.Add(key, "test_data")
	data, _ := cache.Get(key)
	if data == nil {
		t.Errorf("Data not saved %s", data)
	}
}

func TestSlowGet(t *testing.T) {
	key := "key_slow"
	fastCache.Add(key, "test_data")

	time.Sleep(1 * time.Second)

	data, _ := fastCache.Get(key)
	if data == nil {
		t.Error("Data not saved", data)
	}
}

func TestRewriting(t *testing.T) {
	key := "key_to_rewrite"
	cache.Add(key, 40)
	cache.Add(key, 42)
	data, _ := cache.Get(key)
	if data != 42 {
		t.Errorf("Data don't rewriting")
	}
}

func TestDelete(t *testing.T) {
	key := "key_to_delete"
	cache.Add(key, 42)
	cache.Delete(key)
	data, err := cache.Get(key)
	if err == nil {
		t.Error("Don't get error while getting deleted data")
	}
	if data != nil {
		t.Error("Error data", data)
	}
}

func TestSaveToDisk(t *testing.T) {
	key := "value_for_save_to_disk"
	cache.Add(key, 42)
	time.Sleep(2 * time.Second)
	_, err := cache.Get(key)
	if err != nil {
		t.Error("Don't save to disk", err)
	}
}

func TestCacheClear(t *testing.T) {
	fastCache.Add("some", 42)
	fastCache.Add("another", "value")

	currentData := fastCache.GetOnlyInMemory()
	if !(len(currentData) > 0) {
		t.Error("Empty data in memory")
	}

	time.Sleep(2 * time.Second)
	afterIntervalData := cache.GetOnlyInMemory()
	if len(afterIntervalData) > 0 {
		t.Error("Data don't clear", afterIntervalData)
	}
}

func TestRemoveFile(t *testing.T) {
	key := "key_to_remove"
	fastCache.Add(key, 7)

	time.Sleep(1 * time.Second)

	fastCache.Delete(key)
	_, err := fastCache.Get(key)
	if err == nil {
		t.Error("File was not removed")
	}
}
