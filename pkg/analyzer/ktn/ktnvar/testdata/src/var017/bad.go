// Bad examples for the var017 test case.
package var017

import (
	"sync"
	"sync/atomic"
)

// counter avec mutex embarqué (problème si copié).
// Cette structure contient un mutex qui ne doit pas être copié.
type counter struct {
	mu    sync.Mutex
	value int
}

// newCounter crée un nouveau counter.
//
// Returns:
//   - *counter: nouveau compteur
func newCounter() *counter {
	// Retour d'une nouvelle instance
	return &counter{}
}

// increment incrémente le compteur (receiver par valeur - copie le mutex).
func (c counter) increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// safeCounter avec RWMutex embarqué.
// Cette structure contient un RWMutex qui ne doit pas être copié.
type safeCounter struct {
	mu    sync.RWMutex
	value int
}

// newSafeCounter crée un nouveau safeCounter.
//
// Returns:
//   - *safeCounter: nouveau compteur sécurisé
func newSafeCounter() *safeCounter {
	// Retour d'une nouvelle instance
	return &safeCounter{}
}

// read retourne la valeur du compteur (receiver par valeur - copie le RWMutex).
//
// Returns:
//   - int: valeur du compteur.
func (c safeCounter) read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Retourne la valeur actuelle
	return c.value
}

// config avec atomic.Value embarqué.
// Cette structure contient une valeur atomique qui ne doit pas être copiée.
type config struct {
	data atomic.Value
}

// newConfig crée un nouveau config.
//
// Returns:
//   - *config: nouvelle configuration
func newConfig() *config {
	// Retour d'une nouvelle instance
	return &config{}
}

// load charge la configuration (receiver par valeur - copie l'atomic.Value).
//
// Returns:
//   - any: configuration chargée.
func (c config) load() any {
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

// container avec plusieurs types atomiques.
// Cette structure contient plusieurs valeurs atomiques qui ne doivent pas être copiées.
type container struct {
	val1 atomic.Int32
	val2 atomic.Int64
	val3 atomic.Uint32
}

// newContainer crée un nouveau container.
//
// Returns:
//   - *container: nouveau conteneur
func newContainer() *container {
	// Retour d'une nouvelle instance
	return &container{}
}

// values retourne toutes les valeurs atomiques (receiver par valeur).
//
// Returns:
//   - int32: première valeur.
//   - int64: deuxième valeur.
//   - uint32: troisième valeur.
func (c container) values() (int32, int64, uint32) {
	// Retourne les trois valeurs atomiques
	return c.val1.Load(), c.val2.Load(), c.val3.Load()
}

// inner avec mutex imbriqué dans un champ.
// Cette structure contient un mutex qui ne doit pas être copié.
type inner struct {
	mu sync.Mutex
}

// outer avec mutex imbriqué.
// Cette structure contient une struct inner avec mutex.
type outer struct {
	inner inner
}

// newOuter crée un nouveau outer.
//
// Returns:
//   - *outer: nouvelle instance
func newOuter() *outer {
	// Retour d'une nouvelle instance
	return &outer{}
}

// doSomething exécute une action avec le mutex imbriqué (receiver par valeur).
func (o outer) doSomething() {
	o.inner.mu.Lock()
	defer o.inner.mu.Unlock()
}

// flag utilise atomic.Bool.
// Cette structure contient un booléen atomique qui ne doit pas être copié.
type flag struct {
	active atomic.Bool
}

// newFlag crée un nouveau flag.
//
// Returns:
//   - *flag: nouveau flag
func newFlag() *flag {
	// Retour d'une nouvelle instance
	return &flag{}
}

// isActive vérifie si le flag est actif (receiver par valeur - copie atomic.Bool).
//
// Returns:
//   - bool: état du flag.
func (f flag) isActive() bool {
	// Retourne l'état actif
	return f.active.Load()
}

// stats utilise atomic.Uint64.
// Cette structure contient un compteur atomique qui ne doit pas être copié.
type stats struct {
	counter atomic.Uint64
}

// newStats crée un nouveau stats.
//
// Returns:
//   - *stats: nouvelles statistiques
func newStats() *stats {
	// Retour d'une nouvelle instance
	return &stats{}
}

// count retourne le compteur actuel (receiver par valeur - copie atomic.Uint64).
//
// Returns:
//   - uint64: valeur du compteur.
func (s stats) count() uint64 {
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

// init utilise les fonctions privées
func init() {
	// Appel de newCounter
	c := newCounter()
	// Appel de increment
	c.increment()
	// Appel de newSafeCounter
	sc := newSafeCounter()
	// Appel de read
	_ = sc.read()
	// Appel de newConfig
	cfg := newConfig()
	// Appel de load
	_ = cfg.load()
	// Appel de newContainer
	cont := newContainer()
	// Appel de values
	_, _, _ = cont.values()
	// Appel de newOuter
	o := newOuter()
	// Appel de doSomething
	o.doSomething()
	// Appel de newFlag
	f := newFlag()
	// Appel de isActive
	_ = f.isActive()
	// Appel de newStats
	st := newStats()
	// Appel de count
	_ = st.count()
	// Appel de badMutexCopy
	badMutexCopy(sync.Mutex{})
	// Appel de badAssignment
	badAssignment()
	// Appel de badRWMutexParam
	badRWMutexParam(sync.RWMutex{})
	// Appel de badAtomicValueParam
	badAtomicValueParam(atomic.Value{})
	// Appel de badRWMutexAssignment
	badRWMutexAssignment()
	// Appel de badAtomicValueAssignment
	badAtomicValueAssignment()
}
