package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	sync.Mutex
	cap   int               // - capacity
	queue List              // - queue
	items map[Key]*listItem // - items
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	cItem := &cacheItem{
		key:   key,
		value: value,
	}

	l.Lock()
	defer l.Unlock()
	v, ok := l.items[key]
	if ok {
		l.queue.Remove(v)
	} else if l.queue.Len() == l.cap {
		l.purgeBack()
	}
	lItem := l.queue.PushFront(cItem)
	l.items[key] = lItem

	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Lock()
	defer l.Unlock()
	v, ok := l.items[key]
	if !ok {
		return nil, false
	}
	l.queue.MoveToFront(v)

	cValue := v.Value.(*cacheItem)
	return cValue.value, true
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*listItem)
}

func (l *lruCache) purgeBack() {
	back := l.queue.Back()
	l.queue.Remove(back)
	cValue := back.Value.(*cacheItem)
	delete(l.items, cValue.key)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		cap:   capacity,
		queue: NewList(),
		items: make(map[Key]*listItem),
	}
}
