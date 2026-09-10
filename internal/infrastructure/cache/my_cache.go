package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	c                 sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.c.Lock()

	defer c.c.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func (c *Cache) Get(key string) (interface{}, bool) {

	c.c.RLock()

	defer c.c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {

		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}

	}

	return item.Value, true
}

func (c *Cache) Delete(key string) error {

	c.c.Lock()

	defer c.c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	Body := make(map[string]Item)
	return &Cache{defaultExpiration: defaultExpiration, cleanupInterval: cleanupInterval, items: Body}
}
