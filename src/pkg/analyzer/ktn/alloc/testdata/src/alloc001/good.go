package alloc001

// Cas corrects - pas d'allocations inutiles

// GoodReuse réutilise un pointeur existant
func GoodReuse(s *string) *string {
	return s
}

// GoodStructLiteral utilise une variable locale
func GoodStructLiteral() {
	type Config struct{ Name string }
	cfg := Config{Name: "test"}
	_ = cfg
}

// GoodNoAllocation ne fait pas d'allocation
func GoodNoAllocation(x int) int {
	return x * 2
}

// GoodSlicePreallocated préalloue avec bonne taille
func GoodSlicePreallocated() []int {
	items := make([]int, 0, 100)
	for i := 0; i < 100; i++ {
		items = append(items, i)
	}
	return items
}

// GoodNewWithStruct - new() avec struct est OK (pas un type référence)
func GoodNewWithStruct() {
	type Person struct{ Name string }
	p := new(Person)
	_ = p
}

// GoodNewWithInt - new() avec type primitif est OK
func GoodNewWithInt() {
	i := new(int)
	_ = i
}
