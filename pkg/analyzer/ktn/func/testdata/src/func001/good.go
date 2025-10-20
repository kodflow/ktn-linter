package func001

// Good: Small function with few statements
func SmallFunction() string {
	x := 1
	y := 2
	z := x + y
	_ = z
	return "result"
}

// Good: Function with exactly 35 statements
func ExactlyThirtyFive() {
	a := 1                                                                                                                                            // 1
	b := 2                                                                                                                                            // 2
	c := 3                                                                                                                                            // 3
	d := 4                                                                                                                                            // 4
	e := 5                                                                                                                                            // 5
	f := 6                                                                                                                                            // 6
	g := 7                                                                                                                                            // 7
	h := 8                                                                                                                                            // 8
	i := 9                                                                                                                                            // 9
	j := 10                                                                                                                                           // 10
	k := 11                                                                                                                                           // 11
	l := 12                                                                                                                                           // 12
	m := 13                                                                                                                                           // 13
	n := 14                                                                                                                                           // 14
	o := 15                                                                                                                                           // 15
	p := 16                                                                                                                                           // 16
	q := 17                                                                                                                                           // 17
	r := 18                                                                                                                                           // 18
	s := 19                                                                                                                                           // 19
	t := 20                                                                                                                                           // 20
	u := 21                                                                                                                                           // 21
	v := 22                                                                                                                                           // 22
	w := 23                                                                                                                                           // 23
	x := 24                                                                                                                                           // 24
	y := 25                                                                                                                                           // 25
	z := 26                                                                                                                                           // 26
	aa := 27                                                                                                                                          // 27
	ab := 28                                                                                                                                          // 28
	ac := 29                                                                                                                                          // 29
	ad := 30                                                                                                                                          // 30
	ae := 31                                                                                                                                          // 31
	af := 32                                                                                                                                          // 32
	ag := 33                                                                                                                                          // 33
	ah := 34                                                                                                                                          // 34
	_ = a + b + c + d + e + f + g + h + i + j + k + l + m + n + o + p + q + r + s + t + u + v + w + x + y + z + aa + ab + ac + ad + ae + af + ag + ah // 35
}

// Good: Many comments don't count
func ManyComments() {
	// This is a comment
	// This is another comment
	// This is yet another comment
	// Comments don't count
	// More comments
	// Even more comments
	// Still more comments
	x := 1 // First statement
	// Comment between statements
	y := 2 // Second statement
	// Final comment
	_ = x + y // Third statement
}

// Good: Test functions are exempt
func TestSomething() {
	// This can be as long as needed
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// Good: Benchmark functions are exempt
func BenchmarkSomething() {
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// Good: main function is exempt
func main() {
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}


// Good: Example function is exempt
func ExampleSomething() {
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// Good: Fuzz function is exempt
func FuzzSomething() {
	for i := 0; i < 100; i++ {
		x := i * 2
		y := i * 3
		z := x + y
		_ = z
	}
}

// Good: Function with empty lines and various whitespace
func WithEmptyLines() {

	x := 1

	y := 2

	z := x + y

	_ = z

}

// Good: Function with block comments spanning multiple lines
func WithBlockComments() {
	/*
		This is a long
		block comment
		that spans
		many lines
	*/
	x := 1
	/*
		Another block comment
	*/
	y := 2
	/* inline block */ z := x + y
	_ = z
}

// Good: Function with mixed comment styles
func WithMixedComments() {
	// Line comment
	x := 1
	/*
		Block comment
	*/
	y := 2
	// Another line comment
	z := x + y /* inline block */
	_ = z
}

// Good: Function with only braces on separate lines
func WithBraces() {
	{
		x := 1
		_ = x
	}
	{
		y := 2
		_ = y
	}
}

// Good: Function with nested block comments
func WithNestedComments() {
	/* Start block
	   Still in block
	   More block
	*/x := 1
	y := 2
	/*
		Another block
		with multiple lines
		still in block
	*/
	z := x + y
	_ = z
}
