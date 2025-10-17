package rules_data_structures

// ✅ GOOD: initialiser avec make
func writeInitializedMap() {
	m := make(map[string]int) // ✅ initialisé
	m["key"] = 42
}

func safeConditionalWrite(cond bool) {
	var m map[string]int
	if cond {
		m = make(map[string]int)
	}
	if m != nil { // ✅ vérification
		m["key"] = 1
	}
}
