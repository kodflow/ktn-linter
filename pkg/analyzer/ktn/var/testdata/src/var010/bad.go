package var010

// Bad examples: empty slice literals instead of nil (violates KTN-VAR-010)

// badInitUsers creates empty slice literal
func badInitUsers() []string {
	// Empty slice literal - should use var users []string
	users := []string{} // want "KTN-VAR-010"
	// Retour du résultat
	return users
}

// badProcessData creates multiple empty slices
func badProcessData() {
	// Multiple empty slice literals
	names := []string{}  // want "KTN-VAR-010"
	scores := []int{}    // want "KTN-VAR-010"
	flags := []bool{}    // want "KTN-VAR-010"
	ratios := []float64{} // want "KTN-VAR-010"

	// Utilisation des slices
	_ = names
	_ = scores
	_ = flags
	_ = ratios
}

// badNestedSlice creates empty nested slice
func badNestedSlice() {
	// Empty nested slice literal
	matrix := [][]int{} // want "KTN-VAR-010"
	// Utilisation de la slice
	_ = matrix
}

// badStructSlice creates empty struct slice
func badStructSlice() {
	type User struct {
		Name string
	}
	// Empty struct slice literal
	users := []User{} // want "KTN-VAR-010"
	// Utilisation de la slice
	_ = users
}

// badInterfaceSlice creates empty interface slice
func badInterfaceSlice() {
	// Empty interface slice literal
	items := []interface{}{} // want "KTN-VAR-010"
	// Utilisation de la slice
	_ = items
}
