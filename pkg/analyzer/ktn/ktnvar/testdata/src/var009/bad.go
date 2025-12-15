// Package var009 contains test cases for KTN rules.
package var009

// Constantes pour les valeurs de test
const (
	TestField1Value   int     = 42
	TestField4Value   float64 = 3.14
	TestAgeValue      int     = 30
	TestBalanceValue  float64 = 100.0
	TestStruct1Field1 int     = 1
	TestConfigField1  int     = 10
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
		Field1: TestField1Value,
		Field2: "test",
		Field3: true,
		Field4: TestField4Value,
	}
	_ = data
}

// badLargeStructValue2 utilise une autre grande structure par valeur.
func badLargeStructValue2() {
	// Variable locale de grande structure
	user := LargeStruct{
		Field1: TestAgeValue,
		Field2: "John",
		Field3: true,
		Field4: TestBalanceValue,
	}
	_ = user
}

// badMultipleVars déclare plusieurs grandes structures.
func badMultipleVars() {
	// Première variable
	a := LargeStruct{Field1: TestStruct1Field1}
	// Deuxième variable
	b := LargeStruct{Field2: "test"}
	_, _ = a, b
}

// badVarDecl déclare une grande structure avec var.
func badVarDecl() {
	// Déclaration var explicite
	var config LargeStruct
	config.Field1 = TestConfigField1
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
