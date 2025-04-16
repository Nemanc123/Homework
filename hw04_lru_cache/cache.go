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
	lock     sync.Mutex
}
type output struct {
	number interface{}
	mapKey Key
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	value = output{value, key}
	if l.queue.Len() < l.capacity {
		if _, ok := l.items[key]; ok {
			l.queue.Remove(l.items[key])
			l.items[key] = l.queue.PushFront(value)
			return true
		}
		l.items[key] = l.queue.PushFront(value)
		return false
	}
	if _, ok := l.items[key]; ok {
		l.queue.Remove(l.items[key])
		l.items[key] = l.queue.PushFront(value)
		return true
	}
	temp := l.queue.Back()
	l.queue.Remove(temp)
	delete(l.items, temp.Value.(output).mapKey)
	l.items[key] = l.queue.PushFront(value)
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.items[key]; ok {
		l.queue.MoveToFront(l.items[key])
		return l.items[key].Value.(output).number, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	defer l.lock.Unlock()
	l.lock.Lock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
