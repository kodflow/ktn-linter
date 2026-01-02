// Package var015 contains test cases for KTN rules.
package var015

const (
	// UserIdTwo est un ID d'utilisateur
	UserIdTwo int = 2
	// TestIdOne est un ID de test
	TestIdOne int = 1
	// TestIndexZero est un index de test
	TestIndexZero int = 0
)

// badInitUsers creates a map without capacity.
//
// Returns:
//   - map[string]int: map des utilisateurs
func badInitUsers() map[string]int {
	// Map without capacity hint - VIOLATES VAR-009
	users := make(map[string]int)
	users["alice"] = TestIdOne
	users["bob"] = UserIdTwo
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
	data[TestIdOne] = "test"
	cache["key"] = true
	index[TestIndexZero] = []string{"a", "b"}
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
