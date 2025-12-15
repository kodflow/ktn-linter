// Package comment005 contains test cases for KTN rules.
package comment005

// BadUser documentation insuffisante (1 seule ligne) - VIOLATION
type BadUser struct { // want "KTN-COMMENT-005"
	Name string
	Age  int
}
