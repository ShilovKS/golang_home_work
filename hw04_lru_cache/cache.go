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
	mu       sync.RWMutex
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if item, found := l.items[key]; found {
		item.Value = []interface{}{key, value}
		l.queue.MoveToFront(item)
		return true
	}

	newItem := l.queue.PushFront([]interface{}{key, value})
	l.items[key] = newItem

	if l.queue.Len() > l.capacity {
		lastItem := l.queue.Back()
		if lastItem != nil {
			// Извлекаем ключ из значения элемента
			lastValue := lastItem.Value.([]interface{})
			lastKey := lastValue[0].(Key)
			// Удаляем элемент из мапы
			delete(l.items, lastKey)
			// Удаляем элемент из списка
			l.queue.Remove(lastItem)
		}
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if item, found := l.items[key]; found {
		l.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.items = make(map[Key]*ListItem, l.capacity)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
