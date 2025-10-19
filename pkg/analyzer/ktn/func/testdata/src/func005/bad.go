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

// Bad: Switch with default and many cases
// Complexity = 12 (1 base + 10 cases + 1 if)
func SwitchManyWithDefault(x int) int { // want "KTN-FUNC-005"
	result := 0
	switch x {
	case 1:
		result = 1
	case 2:
		result = 2
	case 3:
		result = 3
	case 4:
		result = 4
	case 5:
		result = 5
	case 6:
		result = 6
	case 7:
		result = 7
	case 8:
		result = 8
	case 9:
		result = 9
	case 10:
		result = 10
	default:
		result = 0
	}
	if result > 5 {
		result++
	}
	return result
}

// Bad: Select with default and many comm cases
// Complexity = 11 (1 base + 5 comm cases + 5 ifs)
func SelectManyWithDefault(ch1, ch2, ch3, ch4, ch5 chan int) int { // want "KTN-FUNC-005"
	result := 0
	select {
	case x := <-ch1:
		if x > 0 {
			result = x
		}
	case y := <-ch2:
		if y > 0 {
			result = y
		}
	case z := <-ch3:
		if z > 0 {
			result = z
		}
	case a := <-ch4:
		if a > 0 {
			result = a
		}
	case b := <-ch5:
		if b > 0 {
			result = b
		}
	default:
		result = 0
	}
	if result > 10 {
		result = 10
	}
	return result
}
