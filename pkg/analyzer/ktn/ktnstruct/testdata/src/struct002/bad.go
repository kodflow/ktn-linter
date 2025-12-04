// Bad examples for the struct004 test case.
package struct002

// BadUser documentation insuffisante (1 seule ligne) - VIOLATION
type BadUser struct { // want "KTN-STRUCT-002"
	Name string
	Age  int
}
