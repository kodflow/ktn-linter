package badempty

// Violations avec interface{} et interfaces vides

// Storage interface vide sans méthodes (violation)
type Storage interface{}

// processor utilise interface{} directement (devrait être any ou type spécifique)
type processor interface {
	Process(data interface{}) interface{}
}

// Cache utilise interface{} dans les paramètres
type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{})
}

// badContainer mauvais nommage + interface{}
type badContainer struct {
	items []interface{}
}

// NewBadContainer constructeur qui retourne struct au lieu d'interface
func NewBadContainer() *badContainer {
	return &badContainer{
		items: make([]interface{}, 0),
	}
}
