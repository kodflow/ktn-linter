package var028

func goodDirect(items []int) {
	// Safe in Go 1.22+
	for _, v := range items {
		go processGood(v)
	}
}

func goodDirectIndex(items []int) {
	// Safe in Go 1.22+
	for i, item := range items {
		go func() {
			useGood(i, item)
		}()
	}
}

func goodDifferentVar(items []int) {
	// OK - not same variable
	for _, v := range items {
		x := v + 1
		go processGood(x)
	}
}

func goodDifferentName(items []int) {
	// OK - new variable has different name
	for _, v := range items {
		val := v
		go processGood(val)
	}
}

func goodNonRangeDeclaration(items []int) {
	// OK - not inside a range loop
	x := 10
	y := x
	go processGood(y)
}

func processGood(v int) {}
func useGood(i, item int) {}
