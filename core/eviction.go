package core

var KeysLimit int = 5

func evictFirst() {
	for k := range store {
		delete(store, k)
		return
	}
}

func evict() {
	evictFirst()
}
