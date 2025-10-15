package goodempty

import "fmt"

// CacheImpl implémente l'interface Cache avec types concrets.
type CacheImpl struct {
	data map[string]string
}

// Get récupère une valeur du cache.
//
// Params:
//   - key: clé de la valeur
//
// Returns:
//   - string: valeur associée à la clé
//   - bool: true si la clé existe
func (c *CacheImpl) Get(key string) (string, bool) {
	val, exists := c.data[key]
	return val, exists
}

// Set stocke une valeur dans le cache.
//
// Params:
//   - key: clé de la valeur
//   - value: valeur à stocker
func (c *CacheImpl) Set(key string, value string) {
	c.data[key] = value
}

// NewCache retourne l'interface Cache au lieu de l'implémentation concrète.
//
// Returns:
//   - Cache: nouvelle instance de cache
func NewCache() Cache {
	return &CacheImpl{
		data: make(map[string]string),
	}
}

// ProcessorImpl implémente l'interface Processor.
type ProcessorImpl struct{}

// Process traite les données et retourne le résultat.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat du traitement
func (p *ProcessorImpl) Process(data string) string {
	fmt.Println(data)
	return data
}

// NewProcessor crée un nouveau processeur.
//
// Returns:
//   - Processor: nouvelle instance de processeur
func NewProcessor() Processor {
	return &ProcessorImpl{}
}

// StringContainer implémente Container pour les strings.
type StringContainer struct {
	items []string
}

// Add ajoute un élément au conteneur.
//
// Params:
//   - item: élément à ajouter
func (sc *StringContainer) Add(item string) {
	sc.items = append(sc.items, item)
}

// GetAll retourne tous les éléments.
//
// Returns:
//   - []string: tous les éléments du conteneur
func (sc *StringContainer) GetAll() []string {
	return sc.items
}

// NewStringContainer crée un nouveau conteneur de strings.
//
// Returns:
//   - Container[string]: nouvelle instance de conteneur
func NewStringContainer() Container[string] {
	return &StringContainer{
		items: make([]string, 0),
	}
}
