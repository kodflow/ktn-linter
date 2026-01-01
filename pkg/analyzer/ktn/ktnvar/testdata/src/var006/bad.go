// Package var006 contains test cases for KTN-VAR-006.
// This file contains ALL cases that MUST trigger KTN-VAR-006 errors.
package var006

// NOTE: We can't shadow type names in the same declaration because
// Go will fail to compile. Instead, we shadow function names which
// don't affect type resolution.

// =============================================================================
// SECTION 1: Shadowing built-in functions
// =============================================================================

var (
	// len shadows the built-in len function
	len int = 100 // want "KTN-VAR-006"

	// cap shadows the built-in cap function
	cap int = 200 // want "KTN-VAR-006"

	// append shadows the built-in append function
	append int = 1 // want "KTN-VAR-006"

	// make shadows the built-in make function
	make int = 2 // want "KTN-VAR-006"

	// new shadows the built-in new function
	new int = 3 // want "KTN-VAR-006"

	// panic shadows the built-in panic function
	panic int = 4 // want "KTN-VAR-006"

	// recover shadows the built-in recover function
	recover int = 5 // want "KTN-VAR-006"

	// print shadows the built-in print function
	print int = 6 // want "KTN-VAR-006"

	// println shadows the built-in println function
	println int = 7 // want "KTN-VAR-006"

	// copy shadows the built-in copy function
	copy int = 8 // want "KTN-VAR-006"

	// delete shadows the built-in delete function
	delete int = 9 // want "KTN-VAR-006"

	// close shadows the built-in close function
	close int = 10 // want "KTN-VAR-006"

	// real shadows the built-in real function
	real int = 11 // want "KTN-VAR-006"

	// imag shadows the built-in imag function
	imag int = 12 // want "KTN-VAR-006"

	// complex shadows the built-in complex function
	complex int = 13 // want "KTN-VAR-006"
)
