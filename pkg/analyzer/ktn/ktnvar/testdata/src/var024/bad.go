// Package var024 contains test cases for KTN-VAR-024.
package var024

// badProcess utilise interface{} au lieu de any.
//
// Params:
//   - _data: donnees a traiter (utilise interface{} au lieu de any)
func badProcess(_data interface{}) {} // want "KTN-VAR-024"

var (
	// badX est une variable interface{} au lieu de any.
	badX interface{} // want "KTN-VAR-024"
)

// BadContainer contient un champ interface{}.
// Cette structure illustre l'anti-pattern a eviter.
type BadContainer struct {
	value interface{} // want "KTN-VAR-024"
}

// badReturns retourne interface{} au lieu de any.
//
// Returns:
//   - interface{}: valeur retournee (devrait etre any)
func badReturns() interface{} { // want "KTN-VAR-024"
	// Retourne nil
	return nil
}

// newBadContainer cree un nouveau BadContainer.
//
// Returns:
//   - *BadContainer: nouvelle instance
func newBadContainer() *BadContainer {
	// Retourne une nouvelle instance
	return &BadContainer{}
}

// badInit utilise les variables et fonctions definies.
func badInit() {
	// Appel de badProcess
	badProcess(nil)
	// Utilisation de badX
	_ = badX
	// Appel de badReturns
	_ = badReturns()
	// Creation de BadContainer
	_ = newBadContainer()
}

// init appelle badInit pour eviter KTN-FUNC-004.
func init() {
	badInit()
}
