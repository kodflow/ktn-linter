// Package gospec_declarations démontre les déclarations selon la spec Go.
// Référence: https://go.dev/ref/spec#Declarations_and_scope
package gospec_declarations

import "fmt"

// Spec Go: Constant declarations
// https://go.dev/ref/spec#Constant_declarations

const ValidSingleConst = 42
const ValidTypedConst int = 42
const ValidStringConst string = "hello"

const (
	ValidGroupedConst1 = 1
	ValidGroupedConst2 = 2
	ValidGroupedConst3 = 3
)

// Iota
const (
	ValidIota0 = iota // 0
	ValidIota1        // 1
	ValidIota2        // 2
)

// Spec Go: Variable declarations
// https://go.dev/ref/spec#Variable_declarations

var ValidSingleVar int = 42
var ValidInferredVar = 42
var ValidMultipleVars int

var (
	ValidGroupedVar1 int = 1
	ValidGroupedVar2 int = 2
)

// Multiple variables same type
var ValidX, ValidY, ValidZ int = 1, 2, 3

// Spec Go: Short variable declarations
// https://go.dev/ref/spec#Short_variable_declarations
func ValidShortDeclarations() {
	x := 42
	y, z := "hello", true
	_ = x
	_ = y
	_ = z
}

// Spec Go: Function declarations
// https://go.dev/ref/spec#Function_declarations

func ValidNoParams() {}

func ValidSingleParam(x int) {}

func ValidMultipleParams(x int, y string) {}

func ValidSameTypeParams(x, y, z int) {}

func ValidReturn() int {
	return 42
}

func ValidMultipleReturns() (int, string) {
	return 42, "hello"
}

func ValidNamedReturns() (x int, y string) {
	x = 42
	y = "hello"
	return
}

func ValidVariadicParams(values ...int) {}

// Spec Go: Method declarations
// https://go.dev/ref/spec#Method_declarations

type ValidCounter struct {
	count int
}

func (c *ValidCounter) ValidPointerReceiver() {
	c.count++
}

func (c ValidCounter) ValidValueReceiver() int {
	return c.count
}

// Spec Go: Type declarations
// https://go.dev/ref/spec#Type_declarations

type ValidNewInt int
type ValidNewString string

type ValidAliasInt = int

type ValidStructType struct {
	Field1 int
	Field2 string
}

type ValidInterfaceType interface {
	Method1()
	Method2() int
}

// Spec Go: Blank identifier
// https://go.dev/ref/spec#Blank_identifier

func ValidBlankIdentifier() {
	_ = 42
	x, _ := ValidMultipleReturns()
	_ = x
}

// Spec Go: Scope
// https://go.dev/ref/spec#Declarations_and_scope

func ValidScope() {
	x := 1 // x declared in function scope
	{
		y := 2 // y declared in inner scope
		_ = y
	}
	// y not accessible here
	_ = x
}

// Spec Go: Package-level declarations
var ValidPackageVar = 42

const ValidPackageConst = 100

func ValidPackageFunc() {
	fmt.Println("valid package function")
}

// Spec Go: Exported identifiers
// https://go.dev/ref/spec#Exported_identifiers

// Exported (starts with uppercase)
var ExportedVar = 42
const ExportedConst = 100

func ExportedFunc() {}

type ExportedType struct{}

// Unexported (starts with lowercase)
var unexportedVar = 42
const unexportedConst = 100

func unexportedFunc() {}

type unexportedType struct{}
