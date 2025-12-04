// Bad examples for the var009 test case.
package var009

const (
	// USER_ID_TWO est un ID d'utilisateur
	USER_ID_TWO int = 2
	// TEST_ID_ONE est un ID de test
	TEST_ID_ONE int = 1
	// TEST_INDEX_ZERO est un index de test
	TEST_INDEX_ZERO int = 0
)

// badInitUsers creates a map without capacity.
//
// Returns:
//   - map[string]int: map des utilisateurs
func badInitUsers() map[string]int {
	// Map without capacity hint - VIOLATES VAR-009
	users := make(map[string]int)
	users["alice"] = TEST_ID_ONE
	users["bob"] = USER_ID_TWO
	// Retour du résultat
	return users
}

// badInitConfig creates a map without capacity in var declaration.
func badInitConfig() {
	// Map without capacity hint - VIOLATES VAR-009
	config := make(map[string]string)
	config["host"] = "localhost"
	// Utilisation de la config
	_ = config
}

// badProcessData creates multiple maps without capacity.
func badProcessData() {
	// Multiple maps without capacity hints - VIOLATES VAR-009 (3 times)
	data := make(map[int]string)
	cache := make(map[string]bool)
	index := make(map[int][]string)

	// Utilisation des maps
	data[TEST_ID_ONE] = "test"
	cache["key"] = true
	index[TEST_INDEX_ZERO] = []string{"a", "b"}
}

// badNestedMap creates a map of maps without capacity.
func badNestedMap() {
	// Map without capacity hint - VIOLATES VAR-009
	nested := make(map[string]map[int]string)
	// Inner map also without capacity - VIOLATES VAR-009
	nested["key"] = make(map[int]string)
	// Utilisation de la map
	_ = nested
}

// init utilise les fonctions privées
func init() {
	// Appel de badInitUsers
	badInitUsers()
	// Appel de badInitConfig
	badInitConfig()
	// Appel de badProcessData
	badProcessData()
	// Appel de badNestedMap
	badNestedMap()
}
