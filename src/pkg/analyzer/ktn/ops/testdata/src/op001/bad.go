package op001

func BadDivisionByZero() int {
	// want `\[KTN-OPS-OP-001\] Division par zéro`
	return 10 / 0
}

func BadModuloByZero() int {
	// want `\[KTN-OPS-OP-001\] Division par zéro`
	return 10 % 0
}

func GoodDivisionWithCheck(x, y int) int {
	if y != 0 {
		return x / y
	}
	return 0
}
