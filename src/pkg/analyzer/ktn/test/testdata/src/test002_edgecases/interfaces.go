package test002_edgecases

// Fichier avec seulement des interfaces groupées - devrait être ignoré
// Renommons en interfaces.go pour que la logique spéciale s'applique
// ServiceA defines the interface.
type ServiceA interface {
	MethodA() error
}
// ServiceB defines the interface.

// ServiceB defines the interface.
type ServiceB interface {
	MethodB() string
}
