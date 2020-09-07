package cache

import (
	"fmt"
	"sync"
)

const (
	LastFollowerIDKey   = "last_follower_id"
	TotalFollowerKey    = "total_followers"
	LastFollowerNameKey = "last_follower_name"
)

type Cache struct {
	sync.RWMutex
	data map[string]string
}

func (c *Cache) Set(key, value string) {
	c.Lock()
	c.data[key] = value
	c.Unlock()
}

func (c *Cache) Get(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.data[key]

	if !ok {
		return "", fmt.Errorf("failed to find key %s", key)
	}

	return v, nil
}

func New() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}
