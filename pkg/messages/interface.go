// Interface messages for KTN-INTERFACE rules.
package messages

// registerInterfaceMessages enregistre les messages INTERFACE.
func registerInterfaceMessages() {
	Register(Message{
		Code:  "KTN-INTERFACE-001",
		Short: "interface privée '%s' non utilisée (code mort)",
		Verbose: `PROBLÈME: L'interface privée '%s' n'est utilisée nulle part.

POURQUOI: Une interface non utilisée:
  - Est du code mort
  - Alourdit la maintenance
  - Induit en erreur

ACTIONS:
  - Si prévue pour plus tard → supprimer (git la garde)
  - Si oubliée après refactoring → supprimer
  - Si utilisée par réflexion → annoter avec commentaire`,
	})
}
