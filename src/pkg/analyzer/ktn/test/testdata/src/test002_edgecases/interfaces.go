package test002_edgecases

// Fichier avec seulement des interfaces groupées - devrait être ignoré
// Renommons en interfaces.go pour que la logique spéciale s'applique
type ServiceA interface {
	MethodA() error
}

type ServiceB interface {
	MethodB() string
}
