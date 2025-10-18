package test002 // want `\[KTN_TEST_002\] Fichier 'interfaces_with_const.go' n'a pas de fichier de test correspondant`

// Fichier avec const + interface - ne devrait PAS être ignoré (pas que des interfaces)
// Le const fait que containsOnlyInterfaces002 retourne false

// MaxSize defines the maximum size.
const MaxSize int = 100

// Processor defines the interface.
type Processor interface {
	Do() error
}
