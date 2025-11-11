package func014

// PublicFunction est une fonction publique.
//
// Returns:
//   - string: message
func PublicFunction() string {
	// Appel de la fonction privée
	if processData("test") {
		// Retour du helper
		return privateHelper()
	}
	// Retour vide
	return ""
}

// privateHelper est utilisée par PublicFunction.
//
// Returns:
//   - string: message
func privateHelper() string {
	// Retour du message
	return "helper"
}

// Calculator est une struct.
type Calculator struct {
	value int
}

// Calculate appelle la méthode privée.
//
// Returns:
//   - int: résultat
func (c *Calculator) Calculate() int {
	// Appel de la méthode privée
	return c.compute()
}

// compute est une méthode privée utilisée.
//
// Returns:
//   - int: résultat
func (c *Calculator) compute() int {
	// Retour de la valeur
	return c.value * 2
}

// processData utilise validate en interne.
//
// Params:
//   - data: données
//
// Returns:
//   - bool: succès
func processData(data string) bool {
	// Appel de validate
	return validate(data)
}

// validate est utilisée par processData.
//
// Params:
//   - s: chaîne
//
// Returns:
//   - bool: valide
func validate(s string) bool {
	// Retour de la validation
	return len(s) > 0
}
