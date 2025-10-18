package test002 // want `\[KTN_TEST_002\] Fichier 'interfaces_with_const.go' n'a pas de fichier de test correspondant`

// Fichier avec const + interface - ne devrait PAS être ignoré (pas que des interfaces)
// Le const fait que containsOnlyInterfaces002 retourne false

const MaxSize = 100

type Processor interface {
	Do() error
}
