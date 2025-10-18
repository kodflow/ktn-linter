package map001

func BadMapWriteWithoutCheck(m map[string]int) {
	m["key"] = 42 // want `\[KTN-MAP-001\].*`
}

func BadMapMultipleWrites(m map[string]string) {
	m["a"] = "test" // want `\[KTN-MAP-001\].*`
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
