package var010

// Bad examples: empty slice literals instead of nil (violates KTN-VAR-010)

// badInitUsers creates empty slice literal
func badInitUsers() []string {
	// Empty slice literal - should use var users []string
	users := []string{}
	// Retour du résultat
	return users
}

// badProcessData creates multiple empty slices
func badProcessData() {
	// Multiple empty slice literals
	names := []string{}
	scores := []int{}
	flags := []bool{}
	ratios := []float64{}

	// Utilisation des slices
	_ = names
	_ = scores
	_ = flags
	_ = ratios
}

// badNestedSlice creates empty nested slice
func badNestedSlice() {
	// Empty nested slice literal
	matrix := [][]int{}
	// Utilisation de la slice
	_ = matrix
}

// badStructSlice creates empty struct slice
func badStructSlice() {
	type User struct {
		Name string
	}
	// Empty struct slice literal
	users := []User{}
	// Utilisation de la slice
	_ = users
}

// badInterfaceSlice creates empty interface slice
func badInterfaceSlice() {
	// Empty interface slice literal
	items := []interface{}{}
	// Utilisation de la slice
	_ = items
}
