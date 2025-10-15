// Package rules_var interfaces.
package rules_var

// Validator définit l'interface de validation.
type Validator interface {
	Validate(value interface{}) error
}

// NewValidator crée un nouveau validateur.
//
// Returns:
//   - Validator: l'instance du validateur
func NewValidator() Validator {
	return nil // Placeholder
}
