// Bad examples for the var019 test case.
package var019

import (
	"sync"
	"sync/atomic"
)

// ICounter est l'interface pour Counter.
// Cette interface définit le contrat pour incrémenter le compteur.
type ICounter interface {
	Increment()
}

// Counter avec mutex embarqué (problème si copié).
// Cette structure contient un mutex qui ne doit pas être copié.
type Counter struct {
	mu    sync.Mutex
	value int
}

// NewCounter crée un nouveau Counter.
//
// Returns:
//   - *Counter: nouveau compteur
func NewCounter() *Counter {
	// Retour d'une nouvelle instance
	return &Counter{}
}

// Increment incrémente le compteur (receiver par valeur - copie le mutex).
func (c Counter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// ISafeCounter est l'interface pour SafeCounter.
// Cette interface définit le contrat pour lire le compteur.
type ISafeCounter interface {
	Read() int
}

// SafeCounter avec RWMutex embarqué.
// Cette structure contient un RWMutex qui ne doit pas être copié.
type SafeCounter struct {
	mu    sync.RWMutex
	value int
}

// NewSafeCounter crée un nouveau SafeCounter.
//
// Returns:
//   - *SafeCounter: nouveau compteur sécurisé
func NewSafeCounter() *SafeCounter {
	// Retour d'une nouvelle instance
	return &SafeCounter{}
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

// IConfig est l'interface pour Config.
// Cette interface définit le contrat pour charger la configuration.
type IConfig interface {
	Load() any
}

// Config avec atomic.Value embarqué.
// Cette structure contient une valeur atomique qui ne doit pas être copiée.
type Config struct {
	data atomic.Value
}

// NewConfig crée un nouveau Config.
//
// Returns:
//   - *Config: nouvelle configuration
func NewConfig() *Config {
	// Retour d'une nouvelle instance
	return &Config{}
}

// Load charge la configuration (receiver par valeur - copie l'atomic.Value).
//
// Returns:
//   - any: configuration chargée.
func (c Config) Load() any {
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

// IContainer est l'interface pour Container.
// Cette interface définit le contrat pour accéder aux valeurs atomiques.
type IContainer interface {
	Values() (int32, int64, uint32)
}

// Container avec plusieurs types atomiques.
// Cette structure contient plusieurs valeurs atomiques qui ne doivent pas être copiées.
type Container struct {
	val1 atomic.Int32
	val2 atomic.Int64
	val3 atomic.Uint32
}

// NewContainer crée un nouveau Container.
//
// Returns:
//   - *Container: nouveau conteneur
func NewContainer() *Container {
	// Retour d'une nouvelle instance
	return &Container{}
}

// Values retourne toutes les valeurs atomiques (receiver par valeur).
//
// Returns:
//   - int32: première valeur.
//   - int64: deuxième valeur.
//   - uint32: troisième valeur.
func (c Container) Values() (int32, int64, uint32) {
	// Retourne les trois valeurs atomiques
	return c.val1.Load(), c.val2.Load(), c.val3.Load()
}

// Inner avec mutex imbriqué dans un champ.
// Cette structure contient un mutex qui ne doit pas être copié.
type Inner struct {
	mu sync.Mutex
}

// IOuter est l'interface pour Outer.
// Cette interface définit le contrat pour exécuter une action.
type IOuter interface {
	DoSomething()
}

// Outer avec mutex imbriqué.
// Cette structure contient une struct Inner avec mutex.
type Outer struct {
	inner Inner
}

// NewOuter crée un nouveau Outer.
//
// Returns:
//   - *Outer: nouvelle instance
func NewOuter() *Outer {
	// Retour d'une nouvelle instance
	return &Outer{}
}

// DoSomething exécute une action avec le mutex imbriqué (receiver par valeur).
func (o Outer) DoSomething() {
	o.inner.mu.Lock()
	defer o.inner.mu.Unlock()
}

// IFlag est l'interface pour Flag.
// Cette interface définit le contrat pour vérifier l'état du flag.
type IFlag interface {
	IsActive() bool
}

// Flag utilise atomic.Bool.
// Cette structure contient un booléen atomique qui ne doit pas être copié.
type Flag struct {
	active atomic.Bool
}

// NewFlag crée un nouveau Flag.
//
// Returns:
//   - *Flag: nouveau flag
func NewFlag() *Flag {
	// Retour d'une nouvelle instance
	return &Flag{}
}

// IsActive vérifie si le flag est actif (receiver par valeur - copie atomic.Bool).
//
// Returns:
//   - bool: état du flag.
func (f Flag) IsActive() bool {
	// Retourne l'état actif
	return f.active.Load()
}

// IStats est l'interface pour Stats.
// Cette interface définit le contrat pour accéder au compteur.
type IStats interface {
	Count() uint64
}

// Stats utilise atomic.Uint64.
// Cette structure contient un compteur atomique qui ne doit pas être copié.
type Stats struct {
	counter atomic.Uint64
}

// NewStats crée un nouveau Stats.
//
// Returns:
//   - *Stats: nouvelles statistiques
func NewStats() *Stats {
	// Retour d'une nouvelle instance
	return &Stats{}
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

// init utilise les fonctions privées
func init() {
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
