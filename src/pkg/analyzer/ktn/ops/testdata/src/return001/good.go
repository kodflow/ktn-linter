package return001

// Cas corrects - pas de naked return ou fonctions courtes

func GoodExplicitReturnValue(x int) (result int) {
	result = x * 2
	return result // OK - return explicite
}

func GoodDirectReturn(x int) int {
	return x * 2 // OK - pas de named return
}

func GoodShortNakedReturn(x int) (r int) {
	r = x
	return // OK - fonction courte (< 15 statements)
}

func GoodLongExplicitReturn(v int) (result int) {
	// Fonction longue avec return explicite - OK
	a := v + 1
	b := a * 2
	c := b / 3
	d := c - 4
	e := d + 5
	f := e * 6
	g := f / 7
	h := g - 8
	i := h + 9
	j := i * 10
	k := j / 11
	l := k - 12
	m := l + 13
	n := m * 14
	result = n
	return result // OK - explicit
}

func GoodMultipleReturns(x, y int) (int, int) {
	return x + y, x - y // OK - pas naked
}
