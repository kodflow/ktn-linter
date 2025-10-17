package rules_data_structures

// ✅ GOOD: taille cohérente
func correctArraySize() {
	_ = [3]int{1, 2, 3} // ✅ 3 éléments, taille 3
}

func autoSize() {
	_ = [...]int{1, 2, 3, 4} // ✅ taille automatique
}

// ✅ GOOD: utiliser slice si taille variable
func useSlice() {
	_ = []int{1, 2, 3, 4, 5} // ✅ slice, pas array
}
