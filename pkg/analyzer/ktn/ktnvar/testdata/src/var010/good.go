package var010

// Good examples: using nil slices instead of empty literals (compliant with KTN-VAR-010)

const (
	// VALUE_TWO is constant value 2
	VALUE_TWO int = 2
	// VALUE_THREE is constant value 3
	VALUE_THREE int = 3
	// VALUE_FIVE is constant value 5
	VALUE_FIVE int = 5
	// VALUE_TEN is constant value 10
	VALUE_TEN int = 10
)

// initUsers creates nil slice using var declaration
//
// Returns:
//   - []string: nil slice (no allocation)
func initUsers() []string {
	// Nil slice - good practice (no allocation)
	var users []string
	// Retour du résultat
	return users
}

// processData creates multiple nil slices
func processData() {
	// Multiple nil slices using var declarations
	var names []string
	// scores holds integer scores
	var scores []int
	// flags holds boolean flags
	var flags []bool
	// ratios holds float ratios
	var ratios []float64

	// Utilisation des slices
	_ = names
	_ = scores
	_ = flags
	_ = ratios
}

// nestedSlice creates nil nested slice
func nestedSlice() {
	// Nil nested slice
	var matrix [][]int
	// Utilisation de la slice
	_ = matrix
}

// structSlice creates nil struct slice
func structSlice() {
	type User struct {
		Name string
	}
	// Nil struct slice
	var users []User
	// Utilisation de la slice
	_ = users
}

// interfaceSlice creates nil interface slice
func interfaceSlice() {
	// Nil interface slice
	var items []interface{}
	// Utilisation de la slice
	_ = items
}

// sliceWithElements creates slice with initial elements (allowed)
func sliceWithElements() {
	// Slice with elements - allowed (not empty)
	numbers := []int{1, VALUE_TWO, VALUE_THREE}
	names := []string{"alice", "bob"}

	// Utilisation des slices
	_ = numbers
	_ = names
}

// makeSlice uses make for preallocated slice (allowed)
func makeSlice() {
	// Using make with capacity - allowed
	users := make([]string, 0, VALUE_TEN)
	// Data preallocated with specific length
	var data [VALUE_FIVE]int

	// Utilisation des slices
	_ = users
	_ = data
}

// arrayLiteral uses array literal (not a slice, not checked)
func arrayLiteral() {
	// Array literal - not checked by this rule
	var arr [VALUE_THREE]int
	// Utilisation de l'array
	_ = arr
}
