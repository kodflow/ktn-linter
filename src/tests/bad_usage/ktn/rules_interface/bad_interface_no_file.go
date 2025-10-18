package rules_interface

// ANTI-PATTERN: Ce fichier démontre des structs sans interfaces
// (le package a d'autres fichiers avec interfaces donc pas de violation KTN-INTERFACE-001)

// Ces structs existent sans interfaces dédiées

// SomeStruct quelque struct
type SomeStruct struct {
	data string
}

// DoSomething fait quelque chose
func (s *SomeStruct) DoSomething() error {
	// Early return from function.
	return nil
}

// AnotherStruct encore un struct
type AnotherStruct struct {
	value int
}

// Process traite
func (a *AnotherStruct) Process() {
	// nothing
}
