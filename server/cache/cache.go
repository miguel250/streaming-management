package cache

import (
	"fmt"
	"sync"
	"time"
)

const (
	LastFollowerIDKey   = "last_follower_id"
	TotalFollowerKey    = "total_followers"
	LastFollowerNameKey = "last_follower_name"
	UserAccessCode      = "user_access_code"
	UserAccessExpiresAt = "user_access_expires_at"
	UserRefreshCode     = "user_refresh_code"
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

func (c *Cache) SetAccessToken(token, refreshToken string, expiresIn int64) {
	expires := time.Now().Add(time.Duration(expiresIn) * time.Second)
	c.Set(UserAccessCode, token)
	c.Set(UserAccessExpiresAt, expires.Local().String())
	c.Set(UserRefreshCode, refreshToken)
}

func New() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}
