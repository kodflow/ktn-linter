// Good examples for the var016 test case.
package var015

// Good examples: maps with capacity hints (compliant with KTN-VAR-015)

const (
	// ValueTwo is constant value 2
	ValueTwo int = 2
	// ValueThree is constant value 3
	ValueThree int = 3
	// ValueFive is constant value 5
	ValueFive int = 5
	// ValueTen is constant value 10
	ValueTen int = 10
	// ValueTwenty is constant value 20
	ValueTwenty int = 20
	// ValueFifty is constant value 50
	ValueFifty int = 50
	// ValueHundred is constant value 100
	ValueHundred int = 100
)

// init demonstrates good practices for map initialization with capacity hints
func init() {
	// Map with capacity hint - good practice
	users := make(map[string]int, ValueTen)
	users["alice"] = 1
	users["bob"] = ValueTwo

	// Map with capacity hint in var declaration
	config := make(map[string]string, ValueFive)
	config["host"] = "localhost"
	_ = config

	// Multiple maps with capacity hints
	data := make(map[int]string, ValueHundred)
	cache := make(map[string]bool, ValueFifty)
	index := make(map[int][]string, ValueTwenty)
	data[1] = "test"
	cache["key"] = true
	index[0] = []string{"a", "b"}

	// Map of maps with capacity hint
	nested := make(map[string]map[int]string, ValueTen)
	nested["key"] = make(map[int]string, ValueFive)
	_ = nested

	// Map literal - not checked by this rule
	literal := map[string]int{
		"one":   1,
		"two":   ValueTwo,
		"three": ValueThree,
	}
	_ = literal
}
