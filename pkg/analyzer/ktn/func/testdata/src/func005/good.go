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

// Good: Switch with default (default doesn't add complexity)
// Complexity = 4 (1 base + 3 cases, default = 0)
func SwitchWithDefault(x int) int {
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

// Good: Select with default (default doesn't add complexity)
// Complexity = 3 (1 base + 2 comm cases, default = 0)
func SelectWithDefault(ch1, ch2 chan int) int {
	select {
	case x := <-ch1:
		return x
	case y := <-ch2:
		return y
	default:
		return 0
	}
}

// Good: Function declaration without body (should be skipped)
func ExternalFunc(x int) int

// Good: Range loop (complexity = 2)
func RangeLoop(items []int) int {
	sum := 0
	for _, item := range items {
		sum += item
	}
	return sum
}

// Good: Test function should be skipped (even if complex)
func TestComplexFunction(t int) {
	// This is complex but should be ignored
	if t > 0 {
		t++
	}
	if t > 5 {
		t++
	}
	if t > 10 {
		t++
	}
	for i := 0; i < 3; i++ {
		t++
	}
	for i := 0; i < 3; i++ {
		t++
	}
	for i := 0; i < 3; i++ {
		t++
	}
	switch t {
	case 1:
		t++
	case 2:
		t++
	case 3:
		t++
	case 4:
		t++
	case 5:
		t++
	}
}

// Good: Benchmark function should be skipped (even if complex)
func BenchmarkComplexFunction(b int) {
	// This is complex but should be ignored
	if b > 0 {
		b++
	}
	if b > 5 {
		b++
	}
	if b > 10 {
		b++
	}
	for i := 0; i < 3; i++ {
		b++
	}
	for i := 0; i < 3; i++ {
		b++
	}
	for i := 0; i < 3; i++ {
		b++
	}
	switch b {
	case 1:
		b++
	case 2:
		b++
	case 3:
		b++
	case 4:
		b++
	case 5:
		b++
	}
}

// Good: Example function should be skipped (even if complex)
func ExampleComplexFunction() {
	x := 0
	// This is complex but should be ignored
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
	case 5:
		x++
	}
}

// Good: Fuzz function should be skipped (even if complex)
func FuzzComplexFunction(f int) {
	// This is complex but should be ignored
	if f > 0 {
		f++
	}
	if f > 5 {
		f++
	}
	if f > 10 {
		f++
	}
	for i := 0; i < 3; i++ {
		f++
	}
	for i := 0; i < 3; i++ {
		f++
	}
	for i := 0; i < 3; i++ {
		f++
	}
	switch f {
	case 1:
		f++
	case 2:
		f++
	case 3:
		f++
	case 4:
		f++
	case 5:
		f++
	}
}
