package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mx       *sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l lruCache) Set(key Key, value interface{}) bool {
	l.mx.Lock()
	defer l.mx.Unlock()

	item, ok := l.items[key]
	if !ok {
		newItem := cacheItem{
			key:   key,
			value: value,
		}
		itemList := l.queue.PushFront(newItem)
		l.items[key] = itemList
		if l.queue.Len() > l.capacity {
			lastItem := l.queue.Back()
			l.queue.Remove(lastItem)
			delete(l.items, lastItem.Value.(cacheItem).key)
		}
		return false
	}
	item.Value = cacheItem{
		key:   key,
		value: value,
	}
	l.queue.PushFront(item)
	return true
}

func (l lruCache) Get(key Key) (interface{}, bool) {
	l.mx.Lock()
	defer l.mx.Unlock()

	item, ok := l.items[key]
	if !ok {
		return nil, false
	}
	l.queue.MoveToFront(item)

	return item.Value.(cacheItem).value, true
}

func (l *lruCache) Clear() {
	l.capacity = 0
	l.queue = &list{}
	l.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mx:       &sync.Mutex{},
	}
}
