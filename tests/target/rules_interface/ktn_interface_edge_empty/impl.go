package goodempty

import "fmt"

// CacheImpl implémente l'interface Cache avec types concrets.
type cacheImpl struct {
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
func (c *cacheImpl) Get(key string) (string, bool) {
	val, exists := c.data[key]
	// Retourne la valeur et l'indicateur d'existence
	return val, exists
}

// Set stocke une valeur dans le cache.
//
// Params:
//   - key: clé de la valeur
//   - value: valeur à stocker
func (c *cacheImpl) Set(key string, value string) {
	c.data[key] = value
}

// NewCache retourne l'interface Cache au lieu de l'implémentation concrète.
//
// Returns:
//   - Cache: nouvelle instance de cache
func NewCache() Cache {
	// Retourne une nouvelle instance du cache
	return &cacheImpl{
		data: make(map[string]string),
	}
}

// ProcessorImpl implémente l'interface Processor.
type processorImpl struct{}

// Process traite les données et retourne le résultat.
//
// Params:
//   - data: données à traiter
//
// Returns:
//   - string: résultat du traitement
func (p *processorImpl) Process(data string) string {
	fmt.Println(data)
	// Retourne les données traitées
	return data
}

// NewProcessor crée un nouveau processeur.
//
// Returns:
//   - Processor: nouvelle instance de processeur
func NewProcessor() Processor {
	// Retourne une nouvelle instance du processeur
	return &processorImpl{}
}

// StringContainer implémente Container pour les strings.
type stringContainer struct {
	items []string
}

// Add ajoute un élément au conteneur.
//
// Params:
//   - item: élément à ajouter
func (sc *stringContainer) Add(item string) {
	sc.items = append(sc.items, item)
}

// GetAll retourne tous les éléments.
//
// Returns:
//   - []string: tous les éléments du conteneur
func (sc *stringContainer) GetAll() []string {
	// Retourne tous les éléments du conteneur
	return sc.items
}

// NewStringContainer crée un nouveau conteneur de strings.
//
// Returns:
//   - Container[string]: nouvelle instance de conteneur
func NewStringContainer() Container[string] {
	// Retourne une nouvelle instance du conteneur de strings
	return &stringContainer{
		items: make([]string, 0),
	}
}
