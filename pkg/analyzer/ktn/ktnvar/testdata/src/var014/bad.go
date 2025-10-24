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
	data := LargeStruct{
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
	user := AnotherLargeStruct{
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
	a := LargeStruct{Field1: 1}
	// Deuxième variable
	b := AnotherLargeStruct{Name: "test"}
	_, _ = a, b
}

// badVarDecl déclare une grande structure avec var.
func badVarDecl() {
	// Déclaration var explicite
	var config LargeStruct
	config.Field1 = 10
	_ = config
}
