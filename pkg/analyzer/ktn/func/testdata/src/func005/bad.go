package func005

// Bad: Complexity = 11 (exceeds limit)
func ComplexityEleven(x int) { // want "KTN-FUNC-005"
	if x > 0 {
		x++
	}
	if x > 5 {
		x++
	}
	if x > 10 {
		x++
	}
	for i := 0; i < 3; i++ {
		x++
	}
	for i := 0; i < 3; i++ {
		x++
	}
	for i := 0; i < 3; i++ {
		x++
	}
	switch x {
	case 1:
		x++
	case 2:
		x++
	case 3:
		x++
	case 4:
		x++
	}
}

// Bad: Very complex function
func VeryComplex(x int) int { // want "KTN-FUNC-005"
	if x > 0 {
		if x > 10 {
			if x > 20 {
				return 1
			}
		}
	}

	if x < 0 {
		if x < -10 {
			if x < -20 {
				return -1
			}
		}
	}

	for i := 0; i < x; i++ {
		if i%2 == 0 {
			x++
		} else {
			x--
		}
	}

	switch x {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	default:
		return 0
	}
}

// Bad: Many logical operators
func ManyLogicalOps(a, b, c, d bool) bool { // want "KTN-FUNC-005"
	if a && b {
		return true
	}
	if c || d {
		return true
	}
	if a && c {
		return true
	}
	if b || d {
		return true
	}
	if a && b && c {
		return true
	}
	if c || d || a {
		return true
	}
	if a && b || c && d {
		return true
	}
	return false
}

// Bad: Complex with select and range
func ComplexSelectRange(ch chan int, items []int) int { // want "KTN-FUNC-005"
	result := 0

	// Range adds complexity
	for _, item := range items {
		if item > 0 {
			result += item
		}
	}

	// Select with multiple cases adds complexity
	select {
	case x := <-ch:
		if x > 10 {
			result += x
		}
	case ch <- result:
		if result > 5 {
			result++
		}
	default:
		result = 0
	}

	// More conditions to exceed complexity
	if result > 100 {
		result = 100
	}
	if result < 0 {
		result = 0
	}
	if result%2 == 0 {
		result++
	}
	if result > 50 {
		result = 50
	}

	return result
}
