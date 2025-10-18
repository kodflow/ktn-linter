package rules_var

import "fmt"

// ANTI-PATTERN: Shadowing de variables extrême
// Viole KTN-VAR-003

// Variables du package.
var (
	// GlobalData donnée globale
	GlobalData string = "global"
	// Result résultat global
	Result int = 0
)

// ShadowingNightmare ré-déclare les mêmes noms partout
func ShadowingNightmare() {
	// Shadow de GlobalData au niveau fonction
	GlobalData := "function"
	fmt.Println(GlobalData)

	if true {
		// Shadow encore au niveau if
		GlobalData := "if-block"
		fmt.Println(GlobalData)

		for i := 0; i < 3; i++ {
			// Shadow dans la boucle
			GlobalData := "loop"
			fmt.Println(GlobalData)

			if i == 1 {
				// Shadow encore plus profond
				GlobalData := "nested-if"
				fmt.Println(GlobalData)
			}
		}
	}
}

// ErrorShadowingHell re-déclare err partout
func ErrorShadowingHell() error {
	err := fmt.Errorf("outer")

	if true {
		err := fmt.Errorf("if block") // Shadow !
		fmt.Println(err)

		for i := 0; i < 2; i++ {
			err := fmt.Errorf("loop") // Shadow !
			fmt.Println(err)

			func() {
				err := fmt.Errorf("closure") // Shadow !
				fmt.Println(err)
			}()
		}
	}

	// Return error to caller.
	return err
}

// MultiShadow shadow de plusieurs variables
func MultiShadow() {
	x := 10
	y := 20
	z := 30

	if x > 5 {
		x := 100 // Shadow x
		y := 200 // Shadow y
		z := 300 // Shadow z
		fmt.Println(x, y, z)

		if y > 100 {
			x := 1000 // Shadow encore
			y := 2000 // Shadow encore
			z := 3000 // Shadow encore
			fmt.Println(x, y, z)
		}
	}

	fmt.Println(x, y, z)
}
