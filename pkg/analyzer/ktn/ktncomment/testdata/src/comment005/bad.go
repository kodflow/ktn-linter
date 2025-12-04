// Bad examples for the struct004 test case.
package comment005

// BadUser documentation insuffisante (1 seule ligne) - VIOLATION
type BadUser struct { // want "KTN-COMMENT-005"
	Name string
	Age  int
}
