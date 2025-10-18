package rules_type_ops

// ✅ GOOD: vérification avant déréférencement
func safeDeref() {
	var p *int
	if p != nil { // ✅ vérification
		x := *p
		println(x)
	}
}

func safeDerefWithNew() {
	p := &[]int{0}[0] // ✅ new() retourne non-nil
	*p = 42
	println(*p)
}

func safeCheck(p *int) {
	if p == nil {
		// Early return from function.
		return
	}
	println(*p) // ✅ safe
}
