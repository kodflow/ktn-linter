package var009

// Good examples: maps with capacity hints (compliant with KTN-VAR-009)

const (
	// VALUE_TWO is constant value 2
	VALUE_TWO int = 2
	// VALUE_THREE is constant value 3
	VALUE_THREE int = 3
	// VALUE_FIVE is constant value 5
	VALUE_FIVE int = 5
	// VALUE_TEN is constant value 10
	VALUE_TEN int = 10
	// VALUE_TWENTY is constant value 20
	VALUE_TWENTY int = 20
	// VALUE_FIFTY is constant value 50
	VALUE_FIFTY int = 50
	// VALUE_HUNDRED is constant value 100
	VALUE_HUNDRED int = 100
)

// initUsers creates a map with capacity hint
//
// Returns:
//   - map[string]int: initialized user map with capacity hint
func initUsers() map[string]int {
	// Map with capacity hint - good practice
	users := make(map[string]int, VALUE_TEN)
	users["alice"] = 1
	users["bob"] = VALUE_TWO
	// Retour du résultat
	return users
}

// initConfig creates a map with capacity in var declaration
func initConfig() {
	// Map with capacity hint
	config := make(map[string]string, VALUE_FIVE)
	config["host"] = "localhost"
	// Utilisation de la config
	_ = config
}

// processData creates multiple maps with capacity
func processData() {
	// Multiple maps with capacity hints
	data := make(map[int]string, VALUE_HUNDRED)
	cache := make(map[string]bool, VALUE_FIFTY)
	index := make(map[int][]string, VALUE_TWENTY)

	// Utilisation des maps
	data[1] = "test"
	cache["key"] = true
	index[0] = []string{"a", "b"}
}

// nestedMap creates a map of maps with capacity
func nestedMap() {
	// Map with capacity hint
	nested := make(map[string]map[int]string, VALUE_TEN)
	nested["key"] = make(map[int]string, VALUE_FIVE)
	// Utilisation de la map
	_ = nested
}

// mapLiteral uses map literal (not subject to this rule)
func mapLiteral() {
	// Map literal - not checked by this rule
	data := map[string]int{
		"one":   1,
		"two":   VALUE_TWO,
		"three": VALUE_THREE,
	}
	// Utilisation de la map
	_ = data
}
