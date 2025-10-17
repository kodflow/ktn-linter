package emptyinterfaces

// privateHelper est une interface privée (violation KTN-INTERFACE-007).
// Ce fichier ne devrait contenir QUE des interfaces publiques.
type privateHelper interface {
	help()
}
