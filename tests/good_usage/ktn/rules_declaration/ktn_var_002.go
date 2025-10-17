package rules_declaration

// ✅ GOOD: := pour initialisation
func correctShortDecl() {
	name := "John" // ✅ idiomatique
	age := 30
	active := true
	println(name, age, active)
}

// ✅ GOOD: var OK sans initialisation
func varWithoutInit() {
	var count int // ✅ OK: zero value
	var name string
	println(count, name)
}

// ✅ GOOD: var OK avec type explicite différent
func varWithExplicitType() {
	var x int64 = 42 // ✅ OK: type explicite
	var y float32 = 3.14
	println(x, y)
}
