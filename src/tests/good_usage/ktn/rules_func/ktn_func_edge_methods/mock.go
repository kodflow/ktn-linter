//go:build test
// +build test

package rules_func

// AddConfig configuration pour l'addition multiple.
type AddConfig struct {
	A, B, C, D, E, F float64
}

// MockUserManager est le mock de UserManager.
type MockUserManager struct {
	GetNameFunc  func() string
	SetAgeFunc   func(age int)
	ValidateFunc func() bool
}

// GetName implémente l'interface UserManager.
//
// Returns:
//   - string: le nom de l'utilisateur
func (m *MockUserManager) GetName() string {
	if m.GetNameFunc != nil {
		return m.GetNameFunc()
	}
	return ""
}

// SetAge implémente l'interface UserManager.
//
// Params:
//   - age: l'âge à définir
func (m *MockUserManager) SetAge(age int) {
	if m.SetAgeFunc != nil {
		m.SetAgeFunc(age)
	}
}

// Validate implémente l'interface UserManager.
//
// Returns:
//   - bool: true si valide
func (m *MockUserManager) Validate() bool {
	if m.ValidateFunc != nil {
		return m.ValidateFunc()
	}
	return false
}

// MockCalc est le mock de Calc.
type MockCalc struct {
	AddFunc       func(config AddConfig)
	GetResultFunc func() float64
	CalculateFunc func(x float64) float64
}

// Add implémente l'interface Calc.
//
// Params:
//   - config: configuration de l'addition
func (m *MockCalc) Add(config AddConfig) {
	if m.AddFunc != nil {
		m.AddFunc(config)
	}
}

// GetResult implémente l'interface Calc.
//
// Returns:
//   - float64: le résultat du calcul
func (m *MockCalc) GetResult() float64 {
	if m.GetResultFunc != nil {
		return m.GetResultFunc()
	}
	return 0.0
}

// Calculate implémente l'interface Calc.
//
// Params:
//   - x: la valeur à calculer
//
// Returns:
//   - float64: le résultat du calcul
func (m *MockCalc) Calculate(x float64) float64 {
	if m.CalculateFunc != nil {
		return m.CalculateFunc(x)
	}
	return 0.0
}
