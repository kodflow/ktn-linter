// Bad examples for the struct003 test case.
package struct005

// BadMixedFields représente une struct avec mauvais ordre des champs.
// Les champs exportés sont mélangés avec les champs privés (violation STRUCT-003).
type BadMixedFields struct {
	id        int
	Name      string // want "KTN-STRUCT-005"
	email     string
	Age       int // want "KTN-STRUCT-005"
	visible   bool
	Public    string // want "KTN-STRUCT-005"
	hidden    int
	Exported  bool // want "KTN-STRUCT-005"
	another   string
	LastField int // want "KTN-STRUCT-005"
}

// BadMixedFieldsInterface définit les méthodes de BadMixedFields.
type BadMixedFieldsInterface interface {
	Id() int
	Email() string
	Visible() bool
	Hidden() int
	Another() string
}

// NewBadMixedFields crée une nouvelle instance de BadMixedFields.
//
// Returns:
//   - *BadMixedFields: nouvelle instance
func NewBadMixedFields() *BadMixedFields {
	// Retourne une nouvelle instance
	return &BadMixedFields{}
}

// Id retourne l'identifiant.
//
// Returns:
//   - int: identifiant
func (b *BadMixedFields) Id() int {
	// Retourne le champ id
	return b.id
}

// Email retourne l'email.
//
// Returns:
//   - string: adresse email
func (b *BadMixedFields) Email() string {
	// Retourne le champ email
	return b.email
}

// Visible retourne le statut visible.
//
// Returns:
//   - bool: statut de visibilité
func (b *BadMixedFields) Visible() bool {
	// Retourne le champ visible
	return b.visible
}

// Hidden retourne la valeur cachée.
//
// Returns:
//   - int: valeur cachée
func (b *BadMixedFields) Hidden() int {
	// Retourne le champ hidden
	return b.hidden
}

// Another retourne l'autre valeur.
//
// Returns:
//   - string: autre valeur
func (b *BadMixedFields) Another() string {
	// Retourne le champ another
	return b.another
}
