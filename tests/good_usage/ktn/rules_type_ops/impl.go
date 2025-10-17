package rules_type_ops

// myImplementation implémente MyInterfaceGood.
type myImplementation struct{}

// Do implémente MyInterfaceGood.
func (m *myImplementation) Do() {
	// Implémentation vide pour l'exemple
}

// NewMyInterfaceGood crée une nouvelle instance de MyInterfaceGood.
//
// Returns:
//   - MyInterfaceGood: nouvelle instance
func NewMyInterfaceGood() MyInterfaceGood {
	return &myImplementation{}
}
