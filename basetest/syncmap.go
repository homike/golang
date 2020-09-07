package basetest

import "sync"

var sm sync.Map

func insertKeys() {
	keys := make([]interface{}, 0, 10)
	// Store some keys
	for i := 0; i < 10; i++ {
		v := make([]int, 1)
		keys = append(keys, &v)
		sm.Store(keys[i], struct{}{})
	}
	// delete some keys, but not all keys
	for i, k := range keys {
		if i%2 == 0 {
			continue
		}
		sm.Delete(k)
	}
}

func shutdown() {
	sm.Range(func(key, value interface{}) bool {
		// do something to key
		return true
	})
}

func RunSyncMap() {
	insertKeys()
	// do something ...
	shutdown()
}
