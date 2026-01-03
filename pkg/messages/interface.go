// Package messages provides structured error messages for KTN rules.
// This file contains INTERFACE rule messages.
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

	Register(Message{
		Code:  "KTN-INTERFACE-002",
		Short: "interface '%s' devrait être définie dans le package consommateur",
		Verbose: `PROBLÈME: L'interface '%s' est définie dans le package producteur.

POURQUOI: En Go, les interfaces appartiennent au consommateur:
  - Le consommateur définit ses besoins minimaux
  - Le producteur implémente (sans le savoir)
  - Évite le couplage fort

EFFECTIVE GO:
"Interfaces should be defined by the package that uses them."

EXEMPLE INCORRECT:
  // package producer
  type UserRepository interface { ... }
  type mysqlRepo struct{}

EXEMPLE CORRECT:
  // package consumer (service)
  type UserRepository interface { ... }

  // package producer (repository)
  type MySQLRepository struct{}`,
	})

	Register(Message{
		Code:  "KTN-INTERFACE-003",
		Short: "interface '%s' à une méthode '%s' devrait être nommée '%s'",
		Verbose: `PROBLÈME: Interface '%s' avec méthode '%s' devrait être '%s'.

CONVENTION GO (Effective Go):
  - Interface à une méthode: nom de méthode + "er"
  - Reader, Writer, Closer, Stringer, Handler
  - Rend le code plus lisible et idiomatique

EXEMPLE INCORRECT:
  type Readable interface { Read() }
  type Doable interface { Do() }

EXEMPLE CORRECT:
  type Reader interface { Read() }
  type Doer interface { Do() }`,
	})

	Register(Message{
		Code:  "KTN-INTERFACE-004",
		Short: "interface{} (any) utilisée pour '%s' - préférer interface spécifique",
		Verbose: `PROBLÈME: Utilisation de interface{} (any) pour '%s'.

POURQUOI: interface{} perd la sécurité du typage:
  - Pas de vérification à la compilation
  - Type assertions nécessaires
  - Erreurs runtime possibles

EXCEPTION: Cas légitimes:
  - json.Marshal/Unmarshal
  - Stockage générique (sync.Map, cache)
  - Logging/Debug

EXEMPLE INCORRECT:
  func Process(data interface{}) { ... }

EXEMPLE CORRECT:
  func Process(data Processable) { ... }
  type Processable interface { Process() }`,
	})

	Register(Message{
		Code:  "KTN-INTERFACE-007",
		Short: "paramètre '%s' utilise interface '%s' (%d méthodes), mais n'utilise que %d",
		Verbose: `PROBLÈME: Paramètre '%s' utilise '%s' (%d méthodes) mais n'en utilise que %d.

PRINCIPE DE SÉGRÉGATION DES INTERFACES (ISP):
  "Clients should not be forced to depend on methods they don't use."

POURQUOI:
  - Couplage inutile avec méthodes non utilisées
  - Tests plus difficiles (mocks trop gros)
  - Évolutivité limitée

SOLUTION: Créer une interface minimale avec uniquement les méthodes utilisées.

EXEMPLE INCORRECT:
  func Copy(rw ReadWriter) { rw.Read() } // Write non utilisé

EXEMPLE CORRECT:
  func Copy(r Reader) { r.Read() }`,
	})
}
