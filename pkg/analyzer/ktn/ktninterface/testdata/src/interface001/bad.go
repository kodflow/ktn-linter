// Bad examples for the interface001 test case.
package interface001

// badUnusedInterface is never used anywhere.
type badUnusedInterface interface { // want "KTN-INTERFACE-001"
	Method()
}

// badAnotherUnused is also never referenced.
type badAnotherUnused interface { // want "KTN-INTERFACE-001"
	DoSomething() error
	DoAnother(int) string
}

// badComplexInterface has multiple methods but is unused.
type badComplexInterface interface { // want "KTN-INTERFACE-001"
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
