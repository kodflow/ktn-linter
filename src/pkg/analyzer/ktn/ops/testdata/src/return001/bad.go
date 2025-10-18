package return001

// L'analyzer ne détecte que les naked returns dans les fonctions longues (>10 statements)
func BadNakedReturn(x int) (result int) {
	// Fonction longue pour déclencher l'analyzer
	var a, b, c, d, e int
	a = 1
	b = 2
	c = 3
	d = 4
	e = 5
	result = a + b + c + d + e + x
	result = result * 2
	result = result + 1
	result = result - 1
	return // want `\[KTN-RETURN-001\] Naked return`
}

func BadNakedReturnMultiple(a, b int) (sum int, product int) {
	// Fonction longue
	var x, y, z int
	x = 10
	y = 20
	z = 30
	sum = a + b
	product = a * b
	sum = sum + x
	product = product * y
	sum = sum + z
	product = product / 2
	sum = sum - 5
	product = product + 10
	return // want `\[KTN-RETURN-001\] Naked return`
}

func GoodExplicitReturn(x int) (result int) {
	result = x * 2
	return result
}

func GoodSimpleReturn(x int) int {
	return x * 2
}
