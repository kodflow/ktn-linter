package alloc003

// Cas corrects - slices utilisés correctement

// GoodPreallocatedSlice préalloue avec la bonne taille
func GoodPreallocatedSlice(n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i
	}
	return result
}

// GoodSmallSlice - petite slice sans préallocation est OK
func GoodSmallSlice() []string {
	items := []string{}
	items = append(items, "a")
	items = append(items, "b")
	return items
}

// GoodNoLoop - pas de boucle, pas de problème
func GoodNoLoop() []int {
	return []int{1, 2, 3}
}
