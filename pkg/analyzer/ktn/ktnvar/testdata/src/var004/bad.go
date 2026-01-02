package var004

// Bad examples - names too short

var a int = 1      // want "KTN-VAR-004"
var b string = "x" // want "KTN-VAR-004"

func badExample() {
	x := 42  // want "KTN-VAR-004"
	q := "q" // want "KTN-VAR-004"
	_ = x
	_ = q
}

func badBlockVar() {
	var z int = 10 // want "KTN-VAR-004"
	_ = z
}
