package cache

import (
	"sync"

	"github.com/danillouz/spidey"
)

// Ensure Cache satisfies the CacheService interface.
var _ spidey.CacheService = &Cache{}

// New returns an in-memory Cache which can store visited links.
func New() *Cache {
	return &Cache{
		links: map[string]bool{},
	}
}

// Cache keeps track of visited links and implements a mutex.
type Cache struct {
	links map[string]bool
	mux   sync.Mutex
}

// Add caches a visited link to prevent loops and crawl more efficiently.
func (c *Cache) Add(l string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.links[l]; !ok {
		c.links[l] = true
	}
}

// Check validates if a link has been visited already.
func (c *Cache) Check(l string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok := c.links[l]
	return ok
}
