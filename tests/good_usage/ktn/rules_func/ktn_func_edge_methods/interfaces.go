// Package rules_func contient les fonctions de test edge methods.
package rules_func

// UserManager gère les utilisateurs.
type UserManager interface {
	GetName() string
	SetAge(age int)
	Validate() bool
}

// Calc effectue des calculs.
type Calc interface {
	Add(config AddConfig)
	GetResult() float64
	Calculate(x float64) float64
}
