package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mu       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {

	cache := &lruCache{
		capacity: capacity,
	}
	cache.init()

	return cache
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		item.Value = &cacheItem{key: key, value: value}
		c.queue.MoveToFront(item)
	} else {
		// delete last queue item
		if c.capacity > 0 && c.queue.Len() >= c.capacity {
			tail := c.queue.Back()
			tailValue := c.extractCacheItem(tail)

			c.queue.Remove(tail)
			delete(c.items, tailValue.key)
		}

		c.items[key] = c.queue.PushFront(&cacheItem{key: key, value: value})
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)

	return c.extractCacheItem(item).value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.init()
}

func (c *lruCache) init() {
	c.items = make(map[Key]*listItem, c.capacity)
	c.queue = NewList()
}

func (c *lruCache) extractCacheItem(v *listItem) *cacheItem {
	value, ok := v.Value.(*cacheItem)
	if !ok {
		panic("unexpected value in cache")
	}

	return value
}
