// Bad examples for the struct004 test case.
package struct004

// BadUser documentation insuffisante (1 seule ligne) - VIOLATION
type BadUser struct { // want "KTN-STRUCT-004"
	Name string
	Age  int
}
