package var019

import (
	"sync"
	"sync/atomic"
)

// Counter avec mutex embarqué (problème si copié).
type Counter struct {
	mu    sync.Mutex
	value int
}

// badValueReceiver utilise un receiver par valeur.
func (c Counter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// SafeCounter avec RWMutex embarqué.
type SafeCounter struct {
	mu    sync.RWMutex
	value int
}

// badRWMutexReceiver utilise un receiver par valeur.
func (c SafeCounter) Read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Config avec atomic.Value embarqué.
type Config struct {
	data atomic.Value
}

// badAtomicReceiver utilise un receiver par valeur.
func (c Config) Load() interface{} {
	return c.data.Load()
}

// badMutexCopy copie explicitement un mutex.
func badMutexCopy(mu sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
}

// badAssignment assigne un mutex (copie).
func badAssignment() {
	var mu1 sync.Mutex
	mu2 := mu1
	_ = mu2
}

// Container avec plusieurs types atomiques.
type Container struct {
	val1 atomic.Int32
	val2 atomic.Int64
	val3 atomic.Uint32
}

// badMultipleAtomicFields teste plusieurs champs atomiques.
func (c Container) GetValues() (int32, int64, uint32) {
	return c.val1.Load(), c.val2.Load(), c.val3.Load()
}

// NestedStruct avec mutex imbriqué dans un champ.
type Inner struct {
	mu sync.Mutex
}

type Outer struct {
	inner Inner
}

// badNestedMutex teste un mutex imbriqué.
func (o Outer) DoSomething() {
	o.inner.mu.Lock()
	defer o.inner.mu.Unlock()
}

// badAtomicBool utilise atomic.Bool.
type Flag struct {
	active atomic.Bool
}

func (f Flag) IsActive() bool {
	return f.active.Load()
}

// badAtomicUint64 utilise atomic.Uint64.
type Stats struct {
	counter atomic.Uint64
}

func (s Stats) Count() uint64 {
	return s.counter.Load()
}
