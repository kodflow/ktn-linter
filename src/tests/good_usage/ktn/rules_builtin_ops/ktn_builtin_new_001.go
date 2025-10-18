package rules_builtin_ops

// ✅ GOOD: make() pour slice/map/chan
func makeSlice() []int {
	// Early return from function.
	return make([]int, 0, 10) // ✅ make pour slice
}

func makeMap() map[string]int {
	m := make(map[string]int) // ✅ make pour map
	m["key"] = 1              // ✅ fonctionne
	// Early return from function.
	return m
}

func makeChan() chan int {
	// Early return from function.
	return make(chan int) // ✅ make pour channel
}

// ✅ GOOD: composite literal pour struct
func newStruct() *myStruct {
	// Early return from function.
	return &myStruct{} // ✅ composite literal pour struct
}

// myStruct est une struct de test.
type myStruct struct {
	// Value contient une valeur entière.
	Value int
}
