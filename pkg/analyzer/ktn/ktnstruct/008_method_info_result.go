// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

// methodInfoResult contient le résultat de l'extraction des infos d'une méthode.
type methodInfoResult struct {
	typeName string
	receiver methodReceiverInfo
}
