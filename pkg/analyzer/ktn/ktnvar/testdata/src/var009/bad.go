// Bad examples for the var010 test case.
package var009

// Constantes pour les valeurs de test
const (
	TEST_FIELD1_VALUE   int     = 42
	TEST_FIELD4_VALUE   float64 = 3.14
	TEST_AGE_VALUE      int     = 30
	TEST_BALANCE_VALUE  float64 = 100.0
	TEST_STRUCT1_FIELD1 int     = 1
	TEST_CONFIG_FIELD1  int     = 10
)

// LargeStruct est une structure avec plus de 3 champs.
// Cette structure contient plusieurs champs et est utilisée dans les tests.
type LargeStruct struct {
	Field1 int
	Field2 string
	Field3 bool
	Field4 float64
}

// badLargeStructValue utilise une structure par valeur.
func badLargeStructValue() {
	// Variable locale de grande structure
	data := LargeStruct{
		Field1: TEST_FIELD1_VALUE,
		Field2: "test",
		Field3: true,
		Field4: TEST_FIELD4_VALUE,
	}
	_ = data
}

// badLargeStructValue2 utilise une autre grande structure par valeur.
func badLargeStructValue2() {
	// Variable locale de grande structure
	user := LargeStruct{
		Field1: TEST_AGE_VALUE,
		Field2: "John",
		Field3: true,
		Field4: TEST_BALANCE_VALUE,
	}
	_ = user
}

// badMultipleVars déclare plusieurs grandes structures.
func badMultipleVars() {
	// Première variable
	a := LargeStruct{Field1: TEST_STRUCT1_FIELD1}
	// Deuxième variable
	b := LargeStruct{Field2: "test"}
	_, _ = a, b
}

// badVarDecl déclare une grande structure avec var.
func badVarDecl() {
	// Déclaration var explicite
	var config LargeStruct
	config.Field1 = TEST_CONFIG_FIELD1
	_ = config
}

// init utilise les fonctions privées
func init() {
	// Appel de badLargeStructValue
	badLargeStructValue()
	// Appel de badLargeStructValue2
	badLargeStructValue2()
	// Appel de badMultipleVars
	badMultipleVars()
	// Appel de badVarDecl
	badVarDecl()
}
