package pointer001

// TODO: L'analyzer pointer a besoin d'amélioration pour détecter plus de cas
func GoodNonNilPointer() {
	var p *int
	x := 5
	p = &x
	_ = *p
}
