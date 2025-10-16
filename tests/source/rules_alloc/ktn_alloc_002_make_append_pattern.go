// Package rules_alloc_bad contient des violations KTN-ALLOC-002.
package rules_alloc_bad

// Viole KTN-ALLOC-002 : make([]T, 0) suivi d'append

// BadMakeAppendSimple crée un slice sans capacité puis append.
func BadMakeAppendSimple() {
	items := make([]int, 0) // Viole KTN-ALLOC-002
	items = append(items, 1)
	items = append(items, 2)
	items = append(items, 3)
}

// BadMakeAppendLoop fait des append dans une boucle.
func BadMakeAppendLoop(source []string) {
	result := make([]string, 0) // Viole KTN-ALLOC-002
	for _, v := range source {
		result = append(result, v)
	}
}

// BadMakeAppendRange fait des append avec range.
func BadMakeAppendRange() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	doubled := make([]int, 0) // Viole KTN-ALLOC-002
	for _, n := range numbers {
		doubled = append(doubled, n*2)
	}
}

// BadMakeAppendStruct utilise un slice de struct sans capacité.
func BadMakeAppendStruct() {
	type User struct {
		ID   int
		Name string
	}

	users := make([]User, 0) // Viole KTN-ALLOC-002
	users = append(users, User{ID: 1, Name: "Alice"})
	users = append(users, User{ID: 2, Name: "Bob"})
}

// BadMakeAppendNested fait un append dans un if.
func BadMakeAppendNested(condition bool) {
	items := make([]string, 0) // Viole KTN-ALLOC-002
	if condition {
		items = append(items, "yes")
	} else {
		items = append(items, "no")
	}
}

// BadMakeAppendMultiple fait plusieurs appends d'affilée.
func BadMakeAppendMultiple() {
	data := make([]byte, 0) // Viole KTN-ALLOC-002
	data = append(data, 'h')
	data = append(data, 'e')
	data = append(data, 'l')
	data = append(data, 'l')
	data = append(data, 'o')
}

// Cas farfelus

// BadMakeAppendVariadic utilise append avec variadique.
func BadMakeAppendVariadic() {
	nums := make([]int, 0) // Viole KTN-ALLOC-002
	nums = append(nums, 1, 2, 3, 4, 5)
}

// BadMakeAppendSpread utilise append avec spread.
func BadMakeAppendSpread() {
	a := make([]int, 0) // Viole KTN-ALLOC-002
	b := []int{1, 2, 3}
	a = append(a, b...)
}

// BadMakeAppend2D fait un slice 2D sans capacité.
func BadMakeAppend2D() {
	matrix := make([][]int, 0) // Viole KTN-ALLOC-002
	matrix = append(matrix, []int{1, 2, 3})
	matrix = append(matrix, []int{4, 5, 6})
}

// BadMakeAppendPointers utilise un slice de pointeurs.
func BadMakeAppendPointers() {
	type Node struct {
		Value int
	}
	nodes := make([]*Node, 0) // Viole KTN-ALLOC-002
	nodes = append(nodes, &Node{Value: 1})
	nodes = append(nodes, &Node{Value: 2})
}

// BadMakeAppendInterface utilise un slice d'interfaces.
func BadMakeAppendInterface() {
	items := make([]interface{}, 0) // Viole KTN-ALLOC-002
	items = append(items, 42)
	items = append(items, "hello")
	items = append(items, true)
}

// BadMakeAppendComplexLoop fait un append dans une boucle complexe.
func BadMakeAppendComplexLoop() {
	results := make([]int, 0) // Viole KTN-ALLOC-002
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			results = append(results, i)
		}
	}
}

// BadMakeZeroZeroCapacity utilise make([]T, 0, 0) explicitement.
func BadMakeZeroZeroCapacity() {
	empty := make([]string, 0, 0) // Viole KTN-ALLOC-002
	empty = append(empty, "first")
	empty = append(empty, "second")
}
