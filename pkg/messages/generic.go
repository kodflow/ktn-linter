// Package messages provides structured error messages for KTN rules.
// This file contains GENERIC rule messages.
package messages

// registerGenericMessages enregistre les messages GENERIC.
func registerGenericMessages() {
	Register(Message{
		Code:  "KTN-GENERIC-001",
		Short: "fonction generique '%s' utilise == ou != sans contrainte comparable",
		Verbose: `PROBLEME: La fonction generique '%s' utilise == ou != sur un type parameter avec contrainte 'any'.

POURQUOI: Les operateurs == et != necessitent que le type soit comparable.
La contrainte 'any' accepte des types non-comparables (slices, maps, fonctions).

EXEMPLE INCORRECT:
  func Contains[T any](s []T, v T) bool {
      for _, x := range s {
          if x == v { return true }  // ERREUR: T peut etre non-comparable
      }
      return false
  }

EXEMPLE CORRECT:
  func Contains[T comparable](s []T, v T) bool {
      for _, x := range s {
          if x == v { return true }  // OK: T est comparable
      }
      return false
  }

ALTERNATIVE: Utiliser une fonction de comparaison:
  func Contains[T any](s []T, v T, eq func(T, T) bool) bool {
      for _, x := range s {
          if eq(x, v) { return true }
      }
      return false
  }`,
	})
}
