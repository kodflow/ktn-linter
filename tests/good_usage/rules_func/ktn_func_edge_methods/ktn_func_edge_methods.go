package rules_func

// Méthodes sur types correctement documentées.

// user représente un utilisateur du système.
type user struct {
	name  string
	email string
	age   int
}

// GetName retourne le nom de l'utilisateur.
//
// Returns:
//   - string: le nom de l'utilisateur
func (u *user) GetName() string {
	// Retourne le nom de l'utilisateur
	return u.name
}

// SetAge définit l'âge de l'utilisateur.
//
// Params:
//   - newAge: le nouvel âge à assigner
func (u *user) SetAge(newAge int) {
	u.age = newAge
}

// Validate vérifie que les données utilisateur sont valides.
//
// Returns:
//   - bool: true si l'utilisateur est valide, false sinon
func (u *user) Validate() bool {
	// Vérification du nom non vide
	if u.name == "" {
		// Retourne false car le nom est vide
		return false
	}
	// Vérification de l'email non vide
	if u.email == "" {
		// Retourne false car l'email est vide
		return false
	}
	// Vérification de l'âge dans une plage raisonnable
	if u.age < 0 || u.age > 150 {
		// Retourne false car l'âge est hors limites
		return false
	}
	// Retourne true car tous les champs sont valides
	return true
}

// calculator effectue des calculs et stocke le résultat.
type calculator struct {
	result float64
}

// Add additionne plusieurs nombres (utilise un struct config au lieu de 6+ params).
//
// Params:
//   - config: configuration contenant les valeurs à additionner
func (c *calculator) Add(config AddConfig) {
	c.result = config.A + config.B + config.C + config.D + config.E + config.F
}

// AddConfig configuration pour l'addition multiple.
type AddConfig struct {
	// A est la première valeur à additionner
	A float64
	// B est la deuxième valeur à additionner
	B float64
	// C est la troisième valeur à additionner
	C float64
	// D est la quatrième valeur à additionner
	D float64
	// E est la cinquième valeur à additionner
	E float64
	// F est la sixième valeur à additionner
	F float64
}

// GetResult retourne le résultat du dernier calcul.
//
// Returns:
//   - float64: le résultat stocké
func (c *calculator) GetResult() float64 {
	// Retourne le résultat stocké
	return c.result
}

// Calculate effectue un calcul complexe sur une valeur.
//
// Params:
//   - x: la valeur d'entrée
//
// Returns:
//   - float64: le résultat du calcul
func (c *calculator) Calculate(x float64) float64 {
	// Étape 1: transformation initiale
	result := transformInitial(x)
	// Étape 2: application des opérations
	result = applyOperations(result)
	// Étape 3: normalisation finale
	// Retourne le résultat normalisé
	return normalizeResult(result)
}

// transformInitial transforme la valeur initiale.
//
// Params:
//   - x: valeur à transformer
//
// Returns:
//   - float64: valeur transformée
func transformInitial(x float64) float64 {
	// Retourne la valeur transformée selon la formule
	return (x*2 + 10) * 3
}

// applyOperations applique les opérations.
//
// Params:
//   - x: valeur d'entrée
//
// Returns:
//   - float64: résultat des opérations
func applyOperations(x float64) float64 {
	// Retourne le résultat des opérations appliquées
	return ((x - 5) / 2) + 100
}

// normalizeResult normalise le résultat.
//
// Params:
//   - x: valeur à normaliser
//
// Returns:
//   - float64: valeur normalisée
func normalizeResult(x float64) float64 {
	// Retourne la valeur normalisée
	return (x * 0.5) - 25 + 3.14
}
