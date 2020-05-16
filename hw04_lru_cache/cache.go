package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*cacheItem
}

type cacheItem struct {
	queueItem *listItem
}

func NewCache(capacity int) Cache {

	cache := &lruCache{
		capacity: capacity,
	}
	cache.init()

	return cache
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	val, ok := c.items[key]

	if ok {
		val.queueItem.Value = value
		c.queue.MoveToFront(val.queueItem)
	} else {
		// delete last queue item
		if c.queue.Len() >= c.capacity {
			c.queue.Remove(c.queue.Back())
		}

		c.items[key] = &cacheItem{queueItem: c.queue.PushFront(value)}
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item.queueItem)

	return item.queueItem.Value, true
}

func (c *lruCache) Clear() {
	c.init()
}

func (c *lruCache) init() {
	c.items = make(map[Key]*cacheItem, c.capacity)
	c.queue = NewList()
}