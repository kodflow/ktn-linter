// Package rules_var interfaces.
package rules_var

// Validator définit l'interface de validation.
type Validator interface {
	Validate(value interface{}) error
}
