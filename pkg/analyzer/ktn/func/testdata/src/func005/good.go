package func005

// Good: Simple function (complexity = 1)
func Simple() {
	x := 1
	_ = x
}

// Good: One if statement (complexity = 2)
func OneIf(x int) {
	if x > 0 {
		_ = x
	}
}

// Good: Complexity = 10 (exactly at limit)
// 1 (base) + 3 (if) + 3 (for) + 3 (case) = 10
func ComplexityTen(x int) {
	// +1 for each if
	if x > 0 {
		x++
	}
	if x > 5 {
		x++
	}
	if x > 10 {
		x++
	}

	// +1 for each for
	for i := 0; i < 3; i++ {
		x++
	}
	for i := 0; i < 3; i++ {
		x++
	}
	for i := 0; i < 3; i++ {
		x++
	}

	// +1 for each case
	switch x {
	case 1:
		x++
	case 2:
		x++
	case 3:
		x++
	}
}

// Good: Multiple ifs (complexity = 6)
func MultipleIfs(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	if x == 0 {
		return 0
	}
	if x > 10 {
		return 10
	}
	if x < -10 {
		return -10
	}
	return x
}
