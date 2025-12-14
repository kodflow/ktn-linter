// Package messages provides structured error messages for KTN rules.
// This file contains RETURN rule messages.
package messages

// registerReturnMessages enregistre les messages RETURN.
func registerReturnMessages() {
	Register(Message{
		Code:  "KTN-RETURN-002",
		Short: "retourne nil au lieu de %s vide. Préférer %s{}",
		Verbose: `PROBLÈME: La fonction retourne nil au lieu de %s vide.

POURQUOI: nil peut causer des nil pointer dereferences.
Une collection vide est itérable sans vérification.

EXEMPLE INCORRECT:
  func GetUsers() []User {
      if noUsers {
          return nil  // Danger!
      }
  }

EXEMPLE CORRECT:
  func GetUsers() []User {
      if noUsers {
          return []User{}  // Vide mais safe
      }
  }

NOTE: for range sur nil est OK, mais len() peut surprendre.`,
	})
}
