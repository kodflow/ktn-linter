// Package rules_func contient les fonctions de test edge methods.
package rules_func

// UserManager gère les utilisateurs.
type UserManager interface {
	GetName() string
	SetAge(age int)
	Validate() bool
}

// NewUserManager crée un nouveau gestionnaire d'utilisateurs.
//
// Returns:
//   - UserManager: instance du gestionnaire
func NewUserManager() UserManager {
	return nil // Placeholder
}

// Calc effectue des calculs.
type Calc interface {
	Add(config AddConfig)
	GetResult() float64
	Calculate(x float64) float64
}

// NewCalc crée un nouveau calculateur.
//
// Returns:
//   - Calc: instance du calculateur
func NewCalc() Calc {
	return nil // Placeholder
}
