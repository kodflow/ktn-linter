package func001

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
	*/ x := 1
	y := 2
	/*
	Another block
	with multiple lines
	still in block
	*/
	z := x + y
	_ = z
}
