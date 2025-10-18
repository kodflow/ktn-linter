package interface005_private_only // want `\[KTN_INTERFACE_005\] Fichier interfaces.go existe mais ne contient aucune interface publique`

// Mauvais : interfaces.go existe mais ne contient que des types privés
type privateInterface interface {
	method() error
}

type anotherPrivateInterface interface {
	doSomething()
}
