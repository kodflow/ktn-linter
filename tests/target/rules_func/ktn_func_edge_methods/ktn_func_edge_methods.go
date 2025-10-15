package rules_func

// Méthodes sur types correctement documentées.

// User représente un utilisateur du système.
type User struct {
	name  string
	email string
	age   int
}

// GetName retourne le nom de l'utilisateur.
//
// Returns:
//   - string: le nom de l'utilisateur
func (u *User) GetName() string {
	return u.name
}

// SetAge définit l'âge de l'utilisateur.
//
// Params:
//   - newAge: le nouvel âge à assigner
func (u *User) SetAge(newAge int) {
	u.age = newAge
}

// Validate vérifie que les données utilisateur sont valides.
//
// Returns:
//   - bool: true si l'utilisateur est valide, false sinon
func (u *User) Validate() bool {
	// Vérification du nom non vide
	if u.name == "" {
		return false
	}
	// Vérification de l'email non vide
	if u.email == "" {
		return false
	}
	// Vérification de l'âge dans une plage raisonnable
	if u.age < 0 || u.age > 150 {
		return false
	}
	return true
}

// Calculator effectue des calculs et stocke le résultat.
type Calculator struct {
	result float64
}

// Add additionne plusieurs nombres (utilise un struct config au lieu de 6+ params).
//
// Params:
//   - config: configuration contenant les valeurs à additionner
func (c *Calculator) Add(config AddConfig) {
	c.result = config.A + config.B + config.C + config.D + config.E + config.F
}

// AddConfig configuration pour l'addition multiple.
type AddConfig struct {
	A, B, C, D, E, F float64
}

// GetResult retourne le résultat du dernier calcul.
//
// Returns:
//   - float64: le résultat stocké
func (c *Calculator) GetResult() float64 {
	return c.result
}

// Calculate effectue un calcul complexe sur une valeur.
//
// Params:
//   - x: la valeur d'entrée
//
// Returns:
//   - float64: le résultat du calcul
func (c *Calculator) Calculate(x float64) float64 {
	// Étape 1: transformation initiale
	result := transformInitial(x)
	// Étape 2: application des opérations
	result = applyOperations(result)
	// Étape 3: normalisation finale
	return normalizeResult(result)
}

func transformInitial(x float64) float64 {
	return (x*2 + 10) * 3
}

func applyOperations(x float64) float64 {
	return ((x - 5) / 2) + 100
}

func normalizeResult(x float64) float64 {
	return (x * 0.5) - 25 + 3.14
}
