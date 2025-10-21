package var014

// LargeStruct est une structure avec plus de 3 champs.
type LargeStruct struct {
	Field1 int
	Field2 string
	Field3 bool
	Field4 float64
}

// AnotherLargeStruct est une autre grande structure.
type AnotherLargeStruct struct {
	Name    string
	Age     int
	Email   string
	Active  bool
	Balance float64
}

// badLargeStructValue utilise une structure par valeur.
func badLargeStructValue() {
	// Variable locale de grande structure
	data := LargeStruct{ // want "KTN-VAR-014: utilisez un pointeur pour les structs >64 bytes \\(4 champs détectés\\)"
		Field1: 42,
		Field2: "test",
		Field3: true,
		Field4: 3.14,
	}
	_ = data
}

// badAnotherLargeStructValue utilise une autre grande structure par valeur.
func badAnotherLargeStructValue() {
	// Variable locale de grande structure
	user := AnotherLargeStruct{ // want "KTN-VAR-014: utilisez un pointeur pour les structs >64 bytes \\(5 champs détectés\\)"
		Name:    "John",
		Age:     30,
		Email:   "john@example.com",
		Active:  true,
		Balance: 100.0,
	}
	_ = user
}

// badMultipleVars déclare plusieurs grandes structures.
func badMultipleVars() {
	// Première variable
	a := LargeStruct{Field1: 1} // want "KTN-VAR-014: utilisez un pointeur pour les structs >64 bytes \\(4 champs détectés\\)"
	// Deuxième variable
	b := AnotherLargeStruct{Name: "test"} // want "KTN-VAR-014: utilisez un pointeur pour les structs >64 bytes \\(5 champs détectés\\)"
	_, _ = a, b
}

// badVarDecl déclare une grande structure avec var.
func badVarDecl() {
	// Déclaration var explicite
	var config LargeStruct // want "KTN-VAR-014: utilisez un pointeur pour les structs >64 bytes \\(4 champs détectés\\)"
	config.Field1 = 10
	_ = config
}
