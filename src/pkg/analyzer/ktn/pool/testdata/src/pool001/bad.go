package pool001

import "sync"

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

var myPool = &sync.Pool{
	New: func() interface{} {
		return new(int)
	},
}

func BadPoolGetWithoutDefer() {
	buf := bufferPool.Get().([]byte) // want `\[KTN-POOL-001\].*`
	process(buf)
	// Oubli de retourner au pool
}

func BadPoolGetNoReturn() {
	data := bufferPool.Get().([]byte) // want `\[KTN-POOL-001\].*`
	// Utilisation sans retour
	_ = data
}

func BadPoolGetSimple() {
	item := myPool.Get() // want `\[KTN-POOL-001\].*`
	_ = item
}

func BadPoolGetPointer() {
	val := myPool.Get().(*int) // want `\[KTN-POOL-001\].*`
	*val = 42
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

func GoodPoolPointer() {
	item := myPool.Get()
	defer myPool.Put(item)
	_ = item
}

func process(b []byte) {}
