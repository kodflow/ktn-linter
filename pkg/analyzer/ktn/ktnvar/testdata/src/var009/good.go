// Good examples for the var009 test case.
package var009

const (
	// Answer is the answer
	Answer int = 42
	// AgeValue is age value
	AgeValue int = 25
	// IdValue is id value
	IdValue int = 1
	// UserAge is user age
	UserAge int = 30
	// UserBalance is user balance
	UserBalance float64 = 100.0
	// PiValue is pi value
	PiValue float64 = 3.14
	// ConfigValue is config value
	ConfigValue int = 10
)

// GoodStruct démontre l'utilisation correcte des structures.
// Elle illustre comment allouer des structures avec pointeurs pour optimiser la mémoire.
type GoodStruct struct {
	Field1 int
	Field2 string
	Field3 bool
	Field4 float64
}

// init demonstrates correct usage patterns
func init() {
	// Structure avec pointeur
	data := &GoodStruct{
		Field1: Answer,
		Field2: "test",
		Field3: true,
		Field4: PiValue,
	}
	_ = data

	// Déclaration de pointeur
	var config *GoodStruct
	config = &GoodStruct{Field1: ConfigValue}
	_ = config

	// Allocation avec new
	newData := new(GoodStruct)
	newData.Field1 = Answer
	_ = newData
}
