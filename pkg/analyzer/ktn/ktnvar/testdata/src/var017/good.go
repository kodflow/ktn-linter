// Good examples for the var018 test case.
package var017

import (
	"sync"
	"sync/atomic"
)

const (
	// MULTIPLIER_VALUE multiplicateur pour les calculs
	MULTIPLIER_VALUE int = 2
)

// GoodCounter utilise un pointeur de mutex.
// Compteur thread-safe avec mutex pointeur.
type GoodCounter struct {
	mu    *sync.Mutex // OK: Pointeur de mutex
	value int
}

// GoodCounterInterface définit les méthodes de GoodCounter.
type GoodCounterInterface interface {
	Increment()
	Mu() *sync.Mutex
	Value() int
}

// NewGoodCounter crée un nouveau compteur.
//
// Returns:
//   - *GoodCounter: pointeur vers le nouveau compteur initialisé
func NewGoodCounter() *GoodCounter {
	// Retour de la fonction
	return &GoodCounter{
		mu: &sync.Mutex{},
	}
}

// Increment utilise un receiver par pointeur (correct).
func (c *GoodCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// Mu retourne le mutex.
//
// Returns:
//   - *sync.Mutex: mutex du compteur
func (c *GoodCounter) Mu() *sync.Mutex {
	// Retour du mutex
	return c.mu
}

// Value retourne la valeur.
//
// Returns:
//   - int: valeur du compteur
func (c *GoodCounter) Value() int {
	// Retour de la valeur
	return c.value
}

// GoodSafeCounter utilise un mutex embarqué mais des receivers par pointeur.
// Compteur thread-safe avec RWMutex pour lecture/écriture optimisée.
type GoodSafeCounter struct {
	mu    sync.RWMutex // OK si tous les receivers sont des pointeurs
	value int
}

// GoodSafeCounterInterface définit les méthodes de GoodSafeCounter.
type GoodSafeCounterInterface interface {
	Read() int
	Write(v int)
	Mu() *sync.RWMutex
	Value() int
}

// Read utilise un receiver par pointeur (correct).
//
// Returns:
//   - int: valeur actuelle du compteur
func (c *GoodSafeCounter) Read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Retour de la fonction
	return c.value
}

// Write utilise un receiver par pointeur (correct).
//
// Params:
//   - v: nouvelle valeur à écrire
func (c *GoodSafeCounter) Write(v int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = v
}

// Mu retourne le mutex RW.
//
// Returns:
//   - *sync.RWMutex: mutex du compteur
func (c *GoodSafeCounter) Mu() *sync.RWMutex {
	// Retour du mutex
	return &c.mu
}

// Value retourne la valeur directement (sans lock).
//
// Returns:
//   - int: valeur du compteur
func (c *GoodSafeCounter) Value() int {
	// Retour de la valeur
	return c.value
}

// GoodConfig utilise atomic.Value avec receivers par pointeur.
// Configuration thread-safe avec atomic.Value.
type GoodConfig struct {
	data atomic.Value
}

// GoodConfigInterface définit les méthodes de GoodConfig.
type GoodConfigInterface interface {
	Load() any
	Store(v any)
}

// NewGoodSafeCounter crée un nouveau compteur safe.
//
// Returns:
//   - *GoodSafeCounter: nouveau compteur initialisé
func NewGoodSafeCounter() *GoodSafeCounter {
	// Retour nouvelle instance
	return &GoodSafeCounter{}
}

// NewGoodConfig crée une nouvelle configuration.
//
// Returns:
//   - *GoodConfig: nouvelle configuration
func NewGoodConfig() *GoodConfig {
	// Retour nouvelle instance
	return &GoodConfig{}
}

// Load utilise un receiver par pointeur (correct).
//
// Returns:
//   - any: valeur chargée depuis la configuration
func (c *GoodConfig) Load() any {
	// Retour de la fonction
	return c.data.Load()
}

// Store utilise un receiver par pointeur (correct).
//
// Params:
//   - v: valeur à stocker dans la configuration
func (c *GoodConfig) Store(v any) {
	c.data.Store(v)
}

// goodMutexPointer utilise un pointeur de mutex.
//
// Params:
//   - mu: pointeur vers le mutex à verrouiller
func goodMutexPointer(mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
}

// goodNoAssignment n'assigne pas de mutex.
func goodNoAssignment() {
	mu := &sync.Mutex{} // OK: Pointeur
	mu.Lock()
	defer mu.Unlock()
}

// goodStructWithoutMutex est une struct sans mutex.
type goodStructWithoutMutex struct {
	value int
}

// goodStructWithoutMutexInterface définit les méthodes.
type goodStructWithoutMutexInterface interface {
	ValueMethod() int
}

// ValueMethod peut utiliser receiver par valeur (pas de mutex).
//
// Returns:
//   - int: valeur de la structure
func (s goodStructWithoutMutex) ValueMethod() int {
	// Retour de la fonction
	return s.value
}

// goodTypeAlias est un alias de type (pas une struct).
type goodTypeAlias = int

// goodInterfaceType est une interface (pas une struct avec mutex).
type goodInterfaceType interface {
	DoSomething()
}

// goodImplStruct implémente l'interface.
type goodImplStruct struct {
	value int
}

// useGoodInterfaceType utilise l'interface.
//
// Params:
//   - g: interface à utiliser
func useGoodInterfaceType(g goodInterfaceType) {
	// Utilise l'interface
	g.DoSomething()
}

// DoSomething implémente l'interface (receiver par valeur OK).
func (g goodImplStruct) DoSomething() {
	// Implementation
	_ = g.value
}

// goodPtrMutexStruct utilise un pointeur de mutex comme champ.
type goodPtrMutexStruct struct {
	mu *sync.Mutex
}

// goodPtrMutexStructInterface définit les méthodes.
type goodPtrMutexStructInterface interface {
	ProcessWithPtrMutex()
}

// ProcessWithPtrMutex utilise value receiver (OK car mu est un pointeur).
func (g goodPtrMutexStruct) ProcessWithPtrMutex() {
	// Traitement
	if g.mu != nil {
		g.mu.Lock()
		defer g.mu.Unlock()
	}
}

// goodFuncNoParams est une fonction sans paramètres.
func goodFuncNoParams() {
	// Pas de paramètres donc pas de problème
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()
}

// goodPointerTypeAlias est un alias de type pointeur.
type goodPointerTypeAlias *GoodCounter

// goodNonStructType est un type qui n'est pas une struct.
type goodNonStructType int

// Process traite avec un receiver de type non-struct (OK).
//
// Returns:
//   - int: valeur multipliée
func (g goodNonStructType) Process() int {
	// Retour de la fonction
	return int(g) * MULTIPLIER_VALUE
}

// init utilise les fonctions privées
func init() {
	// Appel de goodMutexPointer
	goodMutexPointer(nil)
	// Appel de goodNoAssignment
	goodNoAssignment()
	// Appel de useGoodInterfaceType
	useGoodInterfaceType(nil)
	// Appel de goodFuncNoParams
	goodFuncNoParams()
}
