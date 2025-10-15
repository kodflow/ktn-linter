package emptyinterfaces

// Ce fichier existe mais ne contient AUCUNE interface publique
// Cela déclenche KTN-INTERFACE-007: interfaces.go vide ou sans interfaces publiques

// privateHelper est une interface privée (ne compte pas).
type privateHelper interface {
	help()
}

// anotherPrivate est une autre interface privée.
type anotherPrivate interface {
	process()
}
