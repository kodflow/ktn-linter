package var027

import "fmt"

// badLoopVariable demonstrates a for loop that should use range int.
func badLoopVariable(n int) {
	for i := 0; i < n; i++ { // want "KTN-VAR-027"
		process(i)
	}
}

// badLoopConstant demonstrates a for loop with constant bound.
func badLoopConstant() {
	for i := 0; i < 10; i++ { // want "KTN-VAR-027"
		fmt.Println(i)
	}
}

// process is a helper function for testing.
func process(i int) {
	fmt.Println(i)
}
