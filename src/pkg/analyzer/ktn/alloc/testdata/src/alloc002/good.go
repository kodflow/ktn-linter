package alloc002

// Cas corrects - append utilisé correctement

// GoodAppendWithPrealloc préalloue la capacité
func GoodAppendWithPrealloc() []int {
	items := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		items = append(items, i)
	}
	return items
}

// GoodAppendReassign réassigne le résultat
func GoodAppendReassign() []int {
	var items []int
	for i := 0; i < 5; i++ {
		items = append(items, i)
	}
	return items
}

// GoodSingleAppend - un seul append est OK
func GoodSingleAppend(items []int) []int {
	return append(items, 42)
}
