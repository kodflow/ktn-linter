package var019

import (
	"sync"
	"sync/atomic"
)

// Counter avec mutex embarqué (problème si copié).
type Counter struct {
	mu    sync.Mutex // want "KTN-VAR-019: struct contient sync\\.Mutex, utiliser \\*Counter pour éviter les copies"
	value int
}

// badValueReceiver utilise un receiver par valeur.
func (c Counter) Increment() { // want "KTN-VAR-019: receiver par valeur copie sync\\.Mutex, utiliser \\*Counter"
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// SafeCounter avec RWMutex embarqué.
type SafeCounter struct {
	mu    sync.RWMutex // want "KTN-VAR-019: struct contient sync\\.RWMutex, utiliser \\*SafeCounter pour éviter les copies"
	value int
}

// badRWMutexReceiver utilise un receiver par valeur.
func (c SafeCounter) Read() int { // want "KTN-VAR-019: receiver par valeur copie sync\\.RWMutex, utiliser \\*SafeCounter"
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Config avec atomic.Value embarqué.
type Config struct {
	data atomic.Value // want "KTN-VAR-019: struct contient atomic\\.Value, utiliser \\*Config pour éviter les copies"
}

// badAtomicReceiver utilise un receiver par valeur.
func (c Config) Load() interface{} { // want "KTN-VAR-019: receiver par valeur copie atomic\\.Value, utiliser \\*Config"
	return c.data.Load()
}

// badMutexCopy copie explicitement un mutex.
func badMutexCopy(mu sync.Mutex) { // want "KTN-VAR-019: passage de sync\\.Mutex par valeur, utiliser \\*sync\\.Mutex"
	mu.Lock()
	defer mu.Unlock()
}

// badAssignment assigne un mutex (copie).
func badAssignment() {
	var mu1 sync.Mutex
	mu2 := mu1 // want "KTN-VAR-019: copie de sync\\.Mutex détectée, utiliser un pointeur"
	_ = mu2
}
