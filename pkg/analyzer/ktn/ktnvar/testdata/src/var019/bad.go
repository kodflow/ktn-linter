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

// Increment incrémente le compteur (receiver par valeur - copie le mutex).
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

// Read retourne la valeur du compteur (receiver par valeur - copie le RWMutex).
//
// Returns:
//   - int: valeur du compteur.
func (c SafeCounter) Read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Retourne la valeur actuelle
	return c.value
}

// Config avec atomic.Value embarqué.
type Config struct {
	data atomic.Value
}

// Load charge la configuration (receiver par valeur - copie l'atomic.Value).
//
// Returns:
//   - interface{}: configuration chargée.
func (c Config) Load() interface{} {
	// Retourne la donnée chargée
	return c.data.Load()
}

// badMutexCopy copie explicitement un mutex.
//
// Params:
//   - mu: mutex passé par valeur (copie).
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

// GetValues retourne toutes les valeurs atomiques (receiver par valeur).
//
// Returns:
//   - int32: première valeur.
//   - int64: deuxième valeur.
//   - uint32: troisième valeur.
func (c Container) GetValues() (int32, int64, uint32) {
	// Retourne les trois valeurs atomiques
	return c.val1.Load(), c.val2.Load(), c.val3.Load()
}

// NestedStruct avec mutex imbriqué dans un champ.
type Inner struct {
	mu sync.Mutex
}

type Outer struct {
	inner Inner
}

// DoSomething exécute une action avec le mutex imbriqué (receiver par valeur).
func (o Outer) DoSomething() {
	o.inner.mu.Lock()
	defer o.inner.mu.Unlock()
}

// badAtomicBool utilise atomic.Bool.
type Flag struct {
	active atomic.Bool
}

// IsActive vérifie si le flag est actif (receiver par valeur - copie atomic.Bool).
//
// Returns:
//   - bool: état du flag.
func (f Flag) IsActive() bool {
	// Retourne l'état actif
	return f.active.Load()
}

// badAtomicUint64 utilise atomic.Uint64.
type Stats struct {
	counter atomic.Uint64
}

// Count retourne le compteur actuel (receiver par valeur - copie atomic.Uint64).
//
// Returns:
//   - uint64: valeur du compteur.
func (s Stats) Count() uint64 {
	// Retourne la valeur du compteur
	return s.counter.Load()
}

// badRWMutexParam passe un RWMutex par valeur.
//
// Params:
//   - mu: RWMutex passé par valeur (copie).
func badRWMutexParam(mu sync.RWMutex) {
	mu.Lock()
	defer mu.Unlock()
}

// badAtomicValueParam passe un atomic.Value par valeur.
//
// Params:
//   - v: atomic.Value passé par valeur (copie).
func badAtomicValueParam(v atomic.Value) {
	_ = v.Load()
}

// badRWMutexAssignment assigne un RWMutex (copie).
func badRWMutexAssignment() {
	var mu1 sync.RWMutex
	mu2 := mu1
	_ = mu2
}

// badAtomicValueAssignment assigne un atomic.Value (copie).
func badAtomicValueAssignment() {
	var v1 atomic.Value
	v2 := v1
	_ = v2
}
