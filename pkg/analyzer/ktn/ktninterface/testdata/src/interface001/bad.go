// Package interface001 contains test cases for KTN rules.
package interface001

// badUnusedInterface is never used anywhere.
type badUnusedInterface interface { // want "KTN-INTERFACE-001: interface privée 'badUnusedInterface' non utilisée \\(code mort\\)"
	Method()
}

// badAnotherUnused is also never referenced.
type badAnotherUnused interface { // want "KTN-INTERFACE-001: interface privée 'badAnotherUnused' non utilisée \\(code mort\\)"
	DoSomething() error
	DoAnother(int) string
}

// badComplexInterface has multiple methods but is unused.
type badComplexInterface interface { // want "KTN-INTERFACE-001: interface privée 'badComplexInterface' non utilisée \\(code mort\\)"
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
