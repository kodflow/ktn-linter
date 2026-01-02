// Package var024 contains test cases for KTN-VAR-024.
package var024

// goodProcess utilise any correctement.
//
// Params:
//   - _data: donnees a traiter (utilise any)
func goodProcess(_data any) {
	// Parametre ignore car fonction de test
}

var (
	// goodX est une variable any correcte.
	goodX any
)

// GoodContainer contient un champ any.
// Cette structure illustre l'utilisation correcte de any.
type GoodContainer struct {
	value any
}

// goodReturns retourne any correctement.
//
// Returns:
//   - any: valeur retournee
func goodReturns() any {
	// Retourne nil
	return nil
}

// Reader est une interface nommee (pas une interface vide).
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer est une interface nommee (pas une interface vide).
type Writer interface {
	Write(p []byte) (n int, err error)
}

// newGoodContainer cree un nouveau GoodContainer.
//
// Returns:
//   - *GoodContainer: nouvelle instance
func newGoodContainer() *GoodContainer {
	// Retourne une nouvelle instance
	return &GoodContainer{}
}

// init utilise les variables et fonctions definies.
func init() {
	// Appel de goodProcess
	goodProcess(nil)
	// Utilisation de goodX
	_ = goodX
	// Appel de goodReturns
	_ = goodReturns()
	// Creation de GoodContainer
	_ = newGoodContainer()
}
