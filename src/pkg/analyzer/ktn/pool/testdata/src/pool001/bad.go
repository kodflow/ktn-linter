package pool001

import "sync"

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func BadPoolGetWithoutDefer() {
	buf := bufferPool.Get().([]byte) // want `\[KTN-POOL-001\] Variable 'buf' obtenue via pool\.Get\(\) sans defer pool\.Put\(\)`
	process(buf)
	// Oubli de retourner au pool
}

func BadPoolGetNoReturn() {
	data := bufferPool.Get().([]byte) // want `\[KTN-POOL-001\] Variable 'data' obtenue via pool\.Get\(\) sans defer pool\.Put\(\)`
	// Utilisation sans retour
	_ = data
}

func GoodPoolGetWithDefer() {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)
	process(buf)
}

func GoodPoolPattern() {
	data := bufferPool.Get().([]byte)
	defer bufferPool.Put(data)
	_ = data
}

func process(b []byte) {}
