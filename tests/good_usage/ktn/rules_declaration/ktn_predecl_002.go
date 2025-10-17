package rules_declaration

// ✅ GOOD: noms différents
func noShadowBuiltins() {
	length := 5        // ✅ length, pas len
	message := "hello" // ✅ message, pas string
	count := 42        // ✅ count, pas int
	println(length, message, count)
}

func noShadowError() {
	errMsg := "failed" // ✅ errMsg, pas error
	_ = errMsg
}

// myError est un type d'erreur personnalisé.
type myError struct{}           // ✅ myError, pas error
func (e myError) Error() string { return "" }
