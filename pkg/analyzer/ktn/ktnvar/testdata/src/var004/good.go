package var004

import "strings"

// Good examples - valid name lengths

var count int = 1      // length >= 2 - OK
var name string = "x"  // length >= 2 - OK

func goodExample() {
	// Loop counters allowed
	for i := 0; i < 10; i++ {
		_ = i
	}

	items := []int{1, 2, 3}
	for j := range items {
		_ = j
	}

	// Type hints allowed
	r := strings.NewReader("")
	m := make(map[string]int)
	ok := true

	_, _, _ = r, m, ok
}
