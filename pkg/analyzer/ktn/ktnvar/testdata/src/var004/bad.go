package var004

// Bad examples - names too short

var a int = 1    // want "KTN-VAR-004"
var b string = "x" // want "KTN-VAR-004"

func badExample() {
	c := 42        // want "KTN-VAR-004"
	_ = c          // use the variable
}
