package cache

import (
	"errors"
	"sync"
	"time"
)

type CacheItem struct {
	Value    interface{}
	PuttedAt time.Time
}

type Cache struct {
	Data          map[string]CacheItem
	InMemoryTime  time.Time
	CheckInterval time.Duration
	Mutex         sync.RWMutex
}

func InitNewCache(checkInterval time.Duration, inMemoryTime time.Time) *Cache {
	if int64(checkInterval) < 0 || inMemoryTime.Unix() < 0 {
		panic("Wrong parameters")
	}

	items := make(map[string]CacheItem)
	cache := &Cache{
		Data:          items,
		InMemoryTime:  inMemoryTime,
		CheckInterval: checkInterval,
	}
	cache.startItemCleaner()

	return cache
}

func (c *Cache) Add(key string, item interface{}) {
	cacheItem := CacheItem{
		Value:    item,
		PuttedAt: time.Now(),
	}
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Data[key] = cacheItem
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	value, isExist := c.Data[key]
	if !isExist {
		return nil, errors.New("Key not found")
	}
	return value, nil
}

func (c *Cache) Delete(key string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	_, isExist := c.Data[key]
	if !isExist {
		return errors.New("Key not found")
	}

	delete(c.Data, key)
	return nil
}

func (c *Cache) startItemCleaner() {
	tiker := time.NewTicker(c.CheckInterval)
	go func() {
		for {
			select {
			case <-tiker.C:
				now := time.Now().Unix()
				for k, v := range c.Data {
					if v.PuttedAt.Unix()+c.InMemoryTime.Unix() < now {
						c.Mutex.Lock()
						delete(c.Data, k)
						c.Mutex.Unlock()
					}
				}
			default:
			}
		}
	}()
}
