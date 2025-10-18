package rules_data_structures

// ✅ GOOD: vérification bounds
func safeIndex(s []int, i int) (int, bool) {
	if i < len(s) { // ✅ check bounds
		return s[i], true
	}
	return 0, false
}

func loopSafe(items []string) {
	for _, item := range items { // ✅ range garantit bounds safety
		println(item)
	}
}
