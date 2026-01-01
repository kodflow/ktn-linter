package var004

// Good examples - valid name lengths

var maxSize int = 100           // length >= 2 - OK
var ok bool = true              // idiomatic short - OK
var id string = "123"           // length >= 2 - OK

// Good - loop variables in loops
func goodExample() {
	for i := 0; i < 10; i++ {
		// i is allowed in loop
	}

	items := []int{1, 2, 3}
	for j, v := range items {
		_, _ = j, v
	}
}
