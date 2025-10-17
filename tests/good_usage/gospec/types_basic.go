// Package gospec_types_basic démontre les types de base selon la spec Go.
// Référence: https://go.dev/ref/spec#Types
package gospec_types_basic

// Spec Go: Boolean types
// https://go.dev/ref/spec#Boolean_types
var validBool bool = true
var validBoolFalse bool = false

// Spec Go: Numeric types
// https://go.dev/ref/spec#Numeric_types

// Integer types
var validUint8 uint8 = 255
var validUint16 uint16 = 65535
var validUint32 uint32 = 4294967295
var validUint64 uint64 = 18446744073709551615
var validInt8 int8 = 127
var validInt16 int16 = 32767
var validInt32 int32 = 2147483647
var validInt64 int64 = 9223372036854775807
var validInt int = 42
var validUint uint = 42

// Floating-point types
var validFloat32 float32 = 3.14
var validFloat64 float64 = 3.14159265358979

// Complex types
var validComplex64 complex64 = 1 + 2i
var validComplex128 complex128 = 1.5 + 2.5i

// Spec Go: String types
// https://go.dev/ref/spec#String_types
var validString string = "hello"
var validRawString string = `raw\nstring`
var validEmptyString string = ""

// Spec Go: Array types
// https://go.dev/ref/spec#Array_types
var validArray [5]int = [5]int{1, 2, 3, 4, 5}
var validArray2 [3]string = [3]string{"a", "b", "c"}
var validArrayInferred [3]int = [...]int{1, 2, 3}

// Spec Go: Slice types
// https://go.dev/ref/spec#Slice_types
var validSlice []int = []int{1, 2, 3}
var validSlice2 []string = make([]string, 5)
var validNilSlice []int

// Spec Go: Struct types
// https://go.dev/ref/spec#Struct_types
type ValidPerson struct {
	Name string
	Age  int
}

var validStruct ValidPerson = ValidPerson{Name: "Alice", Age: 30}
var validStructLiteral = ValidPerson{"Bob", 25}

// Spec Go: Pointer types
// https://go.dev/ref/spec#Pointer_types
var validPointer *int
var validPointerToStruct *ValidPerson = &ValidPerson{Name: "Charlie", Age: 35}

// Spec Go: Function types
// https://go.dev/ref/spec#Function_types
type ValidBinaryOp func(int, int) int

var validFunc ValidBinaryOp = func(a, b int) int {
	return a + b
}

// Spec Go: Interface types
// https://go.dev/ref/spec#Interface_types
type ValidReader interface {
	Read(p []byte) (n int, err error)
}

// Empty interface
var validAny interface{} = 42

// Spec Go: Map types
// https://go.dev/ref/spec#Map_types
var validMap map[string]int = map[string]int{"one": 1, "two": 2}
var validMap2 map[int]string = make(map[int]string)

// Spec Go: Channel types
// https://go.dev/ref/spec#Channel_types
var validChan chan int = make(chan int)
var validChanSend chan<- int = make(chan<- int)
var validChanRecv <-chan int = make(<-chan int)

// Spec Go: Type conversions
// https://go.dev/ref/spec#Conversions
func ValidConversions() {
	var i int = 42
	var f float64 = float64(i)
	var u uint = uint(i)
	_ = f
	_ = u
}

// Spec Go: Type assertions
// https://go.dev/ref/spec#Type_assertions
func ValidTypeAssertion(v interface{}) {
	s, ok := v.(string)
	_ = s
	_ = ok
}

// Spec Go: Zero values
// https://go.dev/ref/spec#The_zero_value
func ValidZeroValues() {
	var zeroInt int        // 0
	var zeroBool bool      // false
	var zeroString string  // ""
	var zeroPointer *int   // nil
	var zeroSlice []int    // nil
	var zeroMap map[int]int // nil
	var zeroChan chan int  // nil
	var zeroFunc func()    // nil

	_ = zeroInt
	_ = zeroBool
	_ = zeroString
	_ = zeroPointer
	_ = zeroSlice
	_ = zeroMap
	_ = zeroChan
	_ = zeroFunc
}
