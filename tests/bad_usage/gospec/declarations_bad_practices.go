// Package gospec_bad_declarations montre des déclarations non-idiomatiques mais qui compilent.
// Ces pratiques violent Effective Go et Go Code Review Comments.
// Référence: https://go.dev/doc/effective_go
// Référence: https://github.com/golang/go/wiki/CodeReviewComments
package gospec_bad_declarations

import "fmt"

// ❌ BAD PRACTICE: Declaring constants individually instead of grouping related ones
// Effective Go recommande de grouper les constantes liées
const BadConst1 = 1
const BadConst2 = 2
const BadConst3 = 3

// ❌ BAD PRACTICE: Declaring variables individually instead of grouping
var badVar1 int = 1
var badVar2 int = 2
var badVar3 int = 3

// ❌ BAD PRACTICE: Not using short variable declaration where appropriate
func BadNoShortDecl() {
	var x int = 42      // Devrait être: x := 42
	var s string = "hi" // Devrait être: s := "hi"
	_, _ = x, s
}

// ❌ BAD PRACTICE: Unnecessary explicit types with obvious initialization
func BadExplicitTypes() {
	var i int = 0          // Devrait être: i := 0
	var f float64 = 0.0    // Devrait être: f := 0.0
	var b bool = false     // Devrait être: b := false
	var str string = ""    // Devrait être: str := ""
	_, _, _, _ = i, f, b, str
}

// ❌ BAD PRACTICE: Not using blank identifier for unused imports/variables
// This would cause compilation error, but pattern of ignoring via assignment is bad
func BadIgnoringValues() {
	x := 42
	y := someFunc() // y not used - should use _ if intentional
	_ = x
	someOtherUse(y) // Added to make it compile
}

// ❌ BAD PRACTICE: snake_case naming (not Go convention)
// Go convention: MixedCaps or mixedCaps
var bad_snake_case_var = 42
const bad_snake_const = 100

func bad_snake_function() {}

type bad_snake_type struct{}

// ❌ BAD PRACTICE: Using var when const would be appropriate
var BadShouldBeConst = 42 // Jamais modifié, devrait être const

// ❌ BAD PRACTICE: Not using named returns for clarity in complex functions
func BadNoNamedReturns(a, b int) (int, int, error) {
	// Dans une fonction complexe, les retours nommés améliorent la clarté
	if a < 0 {
		return 0, 0, fmt.Errorf("invalid")
	}
	return a + b, a * b, nil
}

// ❌ BAD PRACTICE: Redundant type in composite literal
func BadRedundantType() {
	// Le type peut être inféré
	m := map[string]int{
		"a": 1,
	}
	s := []int{1, 2, 3}
	// Usage redondant du type dans append
	s = append(s, []int{4, 5}...)
	_ = m
}

// ❌ BAD PRACTICE: Using new() when literal is clearer
func BadUsingNew() {
	// new() alloue et retourne un pointeur, mais literal est plus clair
	p1 := new(int)
	*p1 = 42
	// Mieux: p1 := &int{} ou simplement travailler avec valeurs

	p2 := new(string)
	*p2 = "hello"
	_, _ = p1, p2
}

// ❌ BAD PRACTICE: Not using type inference
func BadNoTypeInference() {
	var x int = int(42)       // Redondant
	var s string = string("") // Redondant
	_, _ = x, s
}

// ❌ BAD PRACTICE: Verbose variable declarations
func BadVerboseDecls() {
	var err error
	err = nil // Inutile, err est déjà nil

	var counter int
	counter = 0 // Inutile, int zero value est 0

	_ = counter
	_ = err
}

// Helper functions to make examples compile
func someFunc() int { return 0 }
func someOtherUse(int) {}
