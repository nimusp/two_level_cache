package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

const dataPath = "data"

type сacheItem struct {
	Value    interface{}
	PuttedAt time.Duration
}

type Cache struct {
	Data          map[string]сacheItem
	InMemoryTime  time.Duration
	CheckInterval time.Duration
	Mutex         sync.RWMutex
	isWithLog     bool
}

// InitNewCache инициализирует экземпляр кэша с параментами:
// checkInterval - интервал контроля,
// inMemoryTime - время удержания объекта в памяти
func InitNewCache(checkInterval time.Duration, inMemoryTime time.Duration, isWithLog bool) *Cache {
	if int64(checkInterval) < 0 || inMemoryTime < 0 {
		panic("Wrong parameters")
	}

	items := make(map[string]сacheItem)
	cache := &Cache{
		Data:          items,
		InMemoryTime:  inMemoryTime,
		CheckInterval: checkInterval,
		isWithLog:     isWithLog,
	}
	cache.startItemCleaner()

	return cache
}

func (c *Cache) Add(key string, item interface{}) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.addToMap(key, item)
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	valueItem, isExist := c.Data[key]
	var value interface{}
	if !isExist {
		value, isExist = c.readFromFile(key)
		if !isExist {
			return nil, errors.New("Key not found")
		}
		c.addToMap(key, value)
	} else {
		value = valueItem.Value
	}

	return value, nil
}

func (c *Cache) Delete(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.removeFromFile(key)
	delete(c.Data, key)
}

func (c *Cache) getOnlyInMemory() map[string]сacheItem {
	return c.Data
}

func (c *Cache) startItemCleaner() {
	go func() {
		ticker := time.NewTicker(c.CheckInterval)
		for range ticker.C {
			now := time.Now().Unix()
			mapToSave := make(map[string]interface{})
			for k, v := range c.Data {
				at := int64(v.PuttedAt)
				inMemoryTime := int64(c.InMemoryTime / time.Second)
				if at+inMemoryTime < now {
					c.Mutex.Lock()
					mapToSave[k] = v.Value
					c.Mutex.Unlock()
				}
			}
			if len(mapToSave) > 0 {
				c.Mutex.Lock()
				if c.isWithLog {
					log.Println("mapToSave len", len(mapToSave))
					log.Printf("currently %d objects in memory: %+v", len(mapToSave), mapToSave)
				}
				for k, v := range mapToSave {
					c.writeToFile(k, v)
					delete(c.Data, k)
					delete(mapToSave, k)
				}
				if c.isWithLog {
					log.Println("All objects are deleted from memory")
				}
				c.Mutex.Unlock()
			}
		}
	}()
}

func (c *Cache) addToMap(key string, item interface{}) {
	cacheItem := сacheItem{
		Value:    item,
		PuttedAt: time.Duration(time.Now().Unix()),
	}
	c.Data[key] = cacheItem
}

func (c *Cache) readFromFile(key string) (interface{}, bool) {
	fileDir := c.getFileDir() + key

	data, err := ioutil.ReadFile(fileDir)
	if err != nil {
		if c.isWithLog {
			log.Println(err)
		}
		return nil, false
	}
	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, false
	}
	if c.isWithLog {
		log.Println("Read from file", result)
	}
	return result, true
}

func (c *Cache) writeToFile(key string, value interface{}) error {
	rawData, err := json.Marshal(value)

	_, err = os.Stat(dataPath)
	isNotExist := os.IsNotExist(err)
	if isNotExist {
		os.Mkdir(dataPath, 0777)
	}

	fileDir := c.getFileDir() + key
	_, err = os.Stat(fileDir)
	isExist := !os.IsNotExist(err)
	if isExist {
		os.Remove(fileDir)
	}
	err = ioutil.WriteFile(fileDir, rawData, 0644)
	if c.isWithLog {
		log.Println("Wrote to file", value)
	}
	return err
}

func (c *Cache) removeFromFile(key string) {
	fileDir := c.getFileDir() + key
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		return
	}
	os.Remove(fileDir)
}

func (c *Cache) getFileDir() string {
	currentDir, _ := os.Getwd()
	return currentDir + string(os.PathSeparator) + dataPath + string(os.PathSeparator)
}
