package repository

import (
	"sync"

	"Image-Hosting/internal/models"
)

type RedisCache struct {
	mu    sync.RWMutex
	cache map[string]models.Image
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		cache: make(map[string]models.Image),
	}
}

func (c *RedisCache) Get(id string) (models.Image, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	img, exists := c.cache[id]
	return img, exists
}

func (c *RedisCache) Set(id string, img models.Image) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[id] = img
}

func (c *RedisCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, id)
}
