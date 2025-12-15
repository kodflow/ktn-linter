// Bad examples for the struct001 test case.
package struct001

// BadGetterMismatch est une struct avec un getter mal nommé.
// Démontre la violation STRUCT-001: getter Value() retourne le champ data.
type BadGetterMismatch struct {
	data string
}

// BadGetterMismatchInterface définit les méthodes de BadGetterMismatch.
type BadGetterMismatchInterface interface {
	Value() string
}

// NewBadGetterMismatch crée une nouvelle instance.
//
// Returns:
//   - *BadGetterMismatch: nouvelle instance
func NewBadGetterMismatch() *BadGetterMismatch {
	// Retourne une nouvelle instance
	return &BadGetterMismatch{}
}

// Value retourne data mais devrait être nommé Data().
//
// Returns:
//   - string: valeur
func (b *BadGetterMismatch) Value() string { // want "KTN-STRUCT-001: getter 'Value\\(\\)' retourne le champ 'data', devrait être nommé 'Data\\(\\)'"
	// Retour des données
	return b.data
}
