// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

// inconsistencyReport contient les informations pour signaler une incohérence.
type inconsistencyReport struct {
	typeName      string
	minorityType  string
	majorityType  string
	exampleMethod string
}
