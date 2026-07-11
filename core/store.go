package core

import (
	"time"
)

var store map[string]*Obj

type Obj struct {
	Value     interface{}
	ExpiresAt int64
}

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMS int64) *Obj {
	var expiresAt int64 = -1
	if durationMS > 0 {
		expiresAt = time.Now().UnixMilli() + durationMS
	}

	return &Obj{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func Put(k string, obj *Obj) {
	store[k] = obj
}

func Get(k string) *Obj {
	return store[k]
}

func Del(keys []string) int64 {
	keys_deleted := 0
	for _, key := range keys {
		if _, ok := store[key]; ok {
			keys_deleted++
		}
		delete(store, key)
	}
	return int64(keys_deleted)
}
