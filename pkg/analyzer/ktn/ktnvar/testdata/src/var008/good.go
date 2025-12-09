// Good examples for the var009 test case.
package var008

const (
	// ValueTwo is constant value 2
	ValueTwo int = 2
	// ValueThree is constant value 3
	ValueThree int = 3
	// ValueFive is constant value 5
	ValueFive int = 5
	// ValueTen is constant value 10
	ValueTen int = 10
)

// init demonstrates good loop allocation patterns
func init() {
	// Déclaration avant la boucle
	data := make([]int, 0, ValueTen)
	// Loop appends values to reused slice
	for i := range ValueTen {
		// Append current iteration value
		data = append(data, i)
	}
	_ = data

	// Déclaration avant la boucle avec capacité
	cache := make(map[string]int, ValueTen)
	// Loop reuses map
	for i := range ValueTen {
		// Store current value
		cache["key"] = i
	}
	_ = cache

	items := []int{1, ValueTwo, ValueThree}
	// Déclaration avant la boucle avec capacité
	buffer := make([]byte, 0, ValueThree)
	// Range loop reuses buffer
	for _, item := range items {
		// Convert and append item
		buffer = append(buffer, byte(item))
	}
	_ = buffer

	// Pas de boucle, allocation OK avec array
	var single [ValueTen]int
	_ = single

	// Déclaration avant la boucle avec array
	var temp [ValueTen]int
	// Outer loop iterates
	for i := range ValueFive {
		// Inner loop modifies temp
		for j := range ValueFive {
			// Store multiplication result
			temp[j] = i * j
		}
	}
	_ = temp
}
