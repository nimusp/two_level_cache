package cache

import (
	"errors"
	"sync"
	"time"
)

type сacheItem struct {
	Value    interface{}
	PuttedAt time.Duration
}

type Cache struct {
	Data          map[string]сacheItem
	InMemoryTime  time.Duration
	CheckInterval time.Duration
	Mutex         sync.RWMutex
}

// InitNewCache инициализирует экземпляр кэша с параментами:
// checkInterval - интервал контроля,
// inMemoryTime - время удержания объекта в памяти
func InitNewCache(checkInterval time.Duration, inMemoryTime time.Duration) *Cache {
	if int64(checkInterval) < 0 || inMemoryTime < 0 {
		panic("Wrong parameters")
	}

	items := make(map[string]сacheItem)
	cache := &Cache{
		Data:          items,
		InMemoryTime:  inMemoryTime,
		CheckInterval: checkInterval,
	}
	cache.startItemCleaner()

	return cache
}

func (c *Cache) Add(key string, item interface{}) {
	cacheItem := сacheItem{
		Value:    item,
		PuttedAt: time.Duration(time.Now().Nanosecond()),
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
	return value.Value, nil
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
	ticker := time.NewTicker(c.CheckInterval)
	go func() {
		for range ticker.C {
			now := time.Duration(time.Now().Nanosecond())
			for k, v := range c.Data {
				if v.PuttedAt+c.InMemoryTime < now {
					c.Mutex.Lock()
					delete(c.Data, k)
					c.Mutex.Unlock()
				}
			}
		}
	}()
}
