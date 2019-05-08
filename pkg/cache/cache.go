package cache

import (
	"sync"
	"time"

	"github.com/ftl/rigproxy/pkg/protocol"
)

type Cache struct {
	m        map[protocol.CommandKey]entry
	mutex    *sync.RWMutex
	lifetime time.Duration
}

type entry struct {
	resp      protocol.Response
	timestamp time.Time
}

func New() *Cache {
	return NewWithLifetime(0)
}

func NewWithLifetime(lifetime time.Duration) *Cache {
	return &Cache{
		m:        make(map[protocol.CommandKey]entry),
		mutex:    new(sync.RWMutex),
		lifetime: lifetime,
	}
}

func (c *Cache) Put(key protocol.CommandKey, resp protocol.Response) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.m[key] = entry{
		resp:      resp,
		timestamp: time.Now(),
	}
}

func (c *Cache) Get(key protocol.CommandKey) (protocol.Response, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	e, ok := c.m[key]
	if !ok {
		return protocol.Response{}, false
	}
	if c.lifetime > 0 && time.Since(e.timestamp) > c.lifetime {
		return protocol.Response{}, false
	}

	return e.resp, true
}

func (c *Cache) Invalidate(key protocol.CommandKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.m, key)
}
