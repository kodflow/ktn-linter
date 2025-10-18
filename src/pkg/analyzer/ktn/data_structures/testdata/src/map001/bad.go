package map001

func BadMapWriteWithoutCheck(m map[string]int) {
	// want `\[KTN-DS-MAP-001\] Écriture dans une map sans vérification nil`
	m["key"] = 42
}

func BadMapMultipleWrites(m map[string]string) {
	// want `\[KTN-DS-MAP-001\] Écriture dans une map sans vérification nil`
	m["a"] = "test"
	m["b"] = "value"
}

func GoodMapWithCheck(m map[string]int) {
	if m != nil {
		m["key"] = 42
	}
}

func GoodMapInitialized() {
	m := make(map[string]int)
	m["key"] = 42
}
