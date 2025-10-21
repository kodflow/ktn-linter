package var009

// Bad examples: maps without capacity hints (violates KTN-VAR-009)

// badInitUsers creates a map without capacity
func badInitUsers() map[string]int {
	// Map without capacity hint
	users := make(map[string]int) // want "KTN-VAR-009"
	users["alice"] = 1
	users["bob"] = 2
	// Retour du résultat
	return users
}

// badInitConfig creates a map without capacity in var declaration
func badInitConfig() {
	// Map without capacity hint
	var config = make(map[string]string) // want "KTN-VAR-009"
	config["host"] = "localhost"
	// Utilisation de la config
	_ = config
}

// badProcessData creates multiple maps without capacity
func badProcessData() {
	// Multiple maps without capacity hints
	data := make(map[int]string)    // want "KTN-VAR-009"
	cache := make(map[string]bool)  // want "KTN-VAR-009"
	index := make(map[int][]string) // want "KTN-VAR-009"

	// Utilisation des maps
	data[1] = "test"
	cache["key"] = true
	index[0] = []string{"a", "b"}
}

// badNestedMap creates a map of maps without capacity
func badNestedMap() {
	// Map without capacity hint
	nested := make(map[string]map[int]string) // want "KTN-VAR-009"
	nested["key"] = make(map[int]string)      // want "KTN-VAR-009"
	// Utilisation de la map
	_ = nested
}
