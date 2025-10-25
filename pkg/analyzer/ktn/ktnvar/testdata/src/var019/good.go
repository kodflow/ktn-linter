package var019

import (
	"sync"
	"sync/atomic"
)

// GoodCounter utilise un pointeur de mutex.
type GoodCounter struct {
	mu    *sync.Mutex // OK: Pointeur de mutex
	value int
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

// GoodSafeCounter utilise un mutex embarqué mais des receivers par pointeur.
type GoodSafeCounter struct {
	mu    sync.RWMutex // OK si tous les receivers sont des pointeurs
	value int
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

// GoodConfig utilise atomic.Value avec receivers par pointeur.
type GoodConfig struct {
	data atomic.Value
}

// Load utilise un receiver par pointeur (correct).
//
// Returns:
//   - interface{}: valeur chargée depuis la configuration
func (c *GoodConfig) Load() interface{} {
	// Retour de la fonction
	return c.data.Load()
}

// Store utilise un receiver par pointeur (correct).
//
// Params:
//   - v: valeur à stocker dans la configuration
func (c *GoodConfig) Store(v interface{}) {
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

// DoSomething implémente l'interface (receiver par valeur OK).
func (g goodImplStruct) DoSomething() {
	// Implementation
	_ = g.value
}

// goodPtrMutexStruct utilise un pointeur de mutex comme champ.
type goodPtrMutexStruct struct {
	mu *sync.Mutex
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
func (g goodNonStructType) Process() int {
	// Retour de la fonction
	return int(g) * 2
}
