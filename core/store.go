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
	if len(store) >= KeysLimit {
		evict()
	}
	store[k] = obj
}

func Get(k string) *Obj {
	obj := store[k]
	if obj != nil {
		if obj.ExpiresAt <= time.Now().UnixMilli() && obj.ExpiresAt != -1 {
			delete(store, k)
			return nil
		}
	}
	return obj
}

func Del(keys []string) int {
	keys_deleted := 0
	for _, key := range keys {
		if _, ok := store[key]; ok {
			keys_deleted++
		}
		delete(store, key)
	}
	return keys_deleted
}
