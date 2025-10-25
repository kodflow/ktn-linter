package var014

const (
	// ANSWER is the answer
	ANSWER int = 42
	// AGE_VALUE is age value
	AGE_VALUE int = 25
	// ID_VALUE is id value
	ID_VALUE int = 1
	// USER_AGE is user age
	USER_AGE int = 30
	// USER_BALANCE is user balance
	USER_BALANCE float64 = 100.0
	// PI_VALUE is pi value
	PI_VALUE float64 = 3.14
	// CONFIG_VALUE is config value
	CONFIG_VALUE int = 10
)

// GoodSmallStruct est une petite structure (≤3 champs).
type GoodSmallStruct struct {
	ID   int
	Name string
	Age  int
}

// GoodLargeStruct est une grande structure (>3 champs).
type GoodLargeStruct struct {
	Field1 int
	Field2 string
	Field3 bool
	Field4 float64
}

// GoodAnotherLargeStruct est une autre grande structure.
type GoodAnotherLargeStruct struct {
	Name    string
	Age     int
	Email   string
	Active  bool
	Balance float64
}

// goodSmallStructValue utilise une petite structure par valeur.
func goodSmallStructValue() {
	// Petite structure, OK par valeur
	data := GoodSmallStruct{
		ID:   ID_VALUE,
		Name: "test",
		Age:  AGE_VALUE,
	}
	_ = data
}

// goodLargeStructPointer utilise un pointeur pour une grande structure.
func goodLargeStructPointer() {
	// Grande structure avec pointeur
	data := &GoodLargeStruct{
		Field1: ANSWER,
		Field2: "test",
		Field3: true,
		Field4: PI_VALUE,
	}
	_ = data
}

// goodAnotherLargeStructPointer utilise un pointeur.
func goodAnotherLargeStructPointer() {
	// Grande structure avec pointeur
	user := &GoodAnotherLargeStruct{
		Name:    "John",
		Age:     USER_AGE,
		Email:   "john@example.com",
		Active:  true,
		Balance: USER_BALANCE,
	}
	_ = user
}

// goodPointerDecl déclare un pointeur avec var.
func goodPointerDecl() {
	// Déclaration de pointeur
	var config *GoodLargeStruct
	config = &GoodLargeStruct{Field1: CONFIG_VALUE}
	_ = config
}

// goodNewAlloc utilise new pour allouer.
func goodNewAlloc() {
	// Allocation avec new
	data := new(GoodLargeStruct)
	data.Field1 = ANSWER
	_ = data
}
