package utils

import (
	"container/list"
	"sync"
)

type Key interface{}

type entry struct {
	key   interface{}
	value interface{}
}

type LRU struct {
	*list.List
	*sync.RWMutex
	cache   map[interface{}]*list.Element
	maxSize uint64
}

func NewLRU(size uint64) *LRU {
	return &LRU{
		list.New(),
		&sync.RWMutex{},
		make(map[interface{}]*list.Element, size),
		size,
	}
}

func (lru *LRU) Get(key interface{}) (interface{}, bool) {
	defer lru.Unlock()
	lru.Lock()
	if ee, ok := lru.cache[key]; ok {
		lru.MoveToFront(ee)
		return ee.Value, true
	}
	return nil, false
}

func (lru *LRU) Set(key interface{}, value interface{}) {
	defer lru.Unlock()
	lru.Lock()

	if ee, ok := lru.cache[key]; ok {
		lru.MoveToFront(ee)
		ee.Value = value
		return
	}
	if uint64(len(lru.cache)) > lru.maxSize {
		lru.removeOld()
	}
	ele := lru.PushFront(&entry{key: key, value: value})
	lru.cache[key] = ele
}

func (lru *LRU) removeOld() {
	ele := lru.Back()
	if ele != nil {
		lru.Remove(ele)
		key := ele.Value.(*entry).key
		delete(lru.cache, key)
	}
}

func (lru *LRU) pushFront(v interface{}) {
	defer lru.Unlock()
	lru.Lock()
	if ee, ok := lru.cache[v]; ok {
		lru.MoveToFront(ee)
		return
	}
	ele := lru.PushFront(v)
	lru.cache[v] = ele
}

func (lru *LRU) pushBack(v interface{}) {
	defer lru.Unlock()
	lru.Lock()
	if ee, ok := lru.cache[v]; ok {
		lru.MoveToBack(ee)
		return
	}
	ele := lru.PushBack(v)
	lru.cache[v] = ele
}

func (lru *LRU) remove(v interface{}) {
	defer lru.Unlock()
	lru.Lock()

	if ele, hit := lru.cache[v]; hit {
		lru.removeElement(ele)
	}
}

func (lru *LRU) removeElement(e *list.Element) {
	lru.Remove(e)
	kv := e.Value
	delete(lru.cache, kv)
}

func (lru *LRU) list(f func(interface{}) bool) {
	defer lru.Unlock()
	lru.Lock()
	toMove := []*list.Element{}
	for i := lru.Front(); i != nil; i = i.Next() {
		exit := f(i.Value)
		toMove = append(toMove, i)
		if exit {
			break
		}
	}

	for _, v := range toMove {
		lru.MoveToBack(v)
	}

}
