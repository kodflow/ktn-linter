// Package messages provides structured error messages for KTN rules.
// This file contains STRUCT rule messages.
package messages

// registerStructMessages enregistre les messages STRUCT.
func registerStructMessages() {
	Register(Message{
		Code:  "KTN-STRUCT-001",
		Short: "getter '%s' devrait être '%s' pour champ '%s'",
		Verbose: `PROBLÈME: Le getter '%s' ne suit pas la convention.

CONVENTION GO:
  - Champ: name
  - Getter: Name() (pas GetName())
  - Setter: SetName(v)

EXEMPLE INCORRECT:
  func (u *User) GetName() string

EXEMPLE CORRECT:
  func (u *User) Name() string
  func (u *User) SetName(name string)`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-002",
		Short: "struct '%s' sans constructeur New%s()",
		Verbose: `PROBLÈME: La struct '%s' n'a pas de constructeur.

POURQUOI: Un constructeur:
  - Initialise correctement
  - Documente les dépendances
  - Permet la validation
  - Évite structs mal initialisées

EXEMPLE INCORRECT:
  type Service struct { repo Repository }
  // &Service{} → repo est nil!

EXEMPLE CORRECT:
  func NewService(repo Repository) *Service {
      return &Service{repo: repo}
  }`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-003",
		Short: "méthode '%s' utilise Get. Convention Go: %s() pas Get%s()",
		Verbose: `PROBLÈME: La méthode '%s' utilise le préfixe 'Get'.

POURQUOI: Convention Go idiomatique:
  - Getter: Name() pas GetName()
  - Setter: SetName(v) (Set requis)

Effective Go: "It's neither idiomatic nor necessary
to put Get into the getter's name."

EXEMPLE INCORRECT:
  func (u *User) GetName() string

EXEMPLE CORRECT:
  func (u *User) Name() string
  func (u *User) SetName(name string)`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-004",
		Short: "fichier avec %d structs. Une struct par fichier",
		Verbose: `PROBLÈME: Le fichier contient %d structs (max 1).

POURQUOI: Une struct par fichier:
  - Facilite navigation (nom fichier = struct)
  - Évite fichiers de 1000+ lignes
  - Simplifie les reviews
  - Principe de responsabilité unique

EXCEPTION: Structs auxiliaires privées très liées.

SOLUTION: Créer un fichier par struct.
  user.go           → type User struct
  user_repository.go → type UserRepository struct`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-005",
		Short: "champs privés avant publics. Exportés d'abord",
		Verbose: `PROBLÈME: Champs privés avant champs publics.

POURQUOI: Les champs exportés = API publique.
Les mettre en premier facilite la lecture.

ORDRE ATTENDU:
  1. Champs exportés (Majuscule)
  2. Champs privés (minuscule)

EXEMPLE INCORRECT:
  type User struct {
      password string  // Privé
      Name     string  // Exporté après!
  }

EXEMPLE CORRECT:
  type User struct {
      Name     string  // Exportés d'abord
      Email    string
      password string  // Puis privés
  }`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-006",
		Short: "champ privé '%s' avec tag de sérialisation inutile",
		Verbose: `PROBLÈME: Le champ privé '%s' a un tag json/xml/yaml.

POURQUOI: Les champs privés ne sont PAS sérialisés.
Le tag est donc inutile et trompeur.

EXEMPLE INCORRECT:
  type User struct {
      name string ` + "`json:\"name\"`" + `  // Ignoré!
  }

EXEMPLE CORRECT (exporter):
  type User struct {
      Name string ` + "`json:\"name\"`" + `
  }

EXEMPLE CORRECT (garder privé):
  type User struct {
      name string  // Pas de tag
  }`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-007",
		Short: "champ exporté '%s' sans tag json/xml dans struct DTO",
		Verbose: `PROBLÈME: Le champ exporté '%s' n'a pas de tag de sérialisation.

POURQUOI: Les structs DTO doivent avoir des tags explicites:
  - Garantit un contrat de sérialisation stable
  - Évite les changements accidentels lors de refactoring
  - Documente le format de l'API

EXEMPLE INCORRECT:
  type UserDTO struct {
      Name  string  // Pas de tag
      Email string
  }

EXEMPLE CORRECT:
  type UserDTO struct {
      Name  string ` + "`json:\"name\"`" + `
      Email string ` + "`json:\"email\"`" + `
  }`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-008",
		Short: "receiver '%s' est %s mais '%s' utilise %s sur type '%s'",
		Verbose: `PROBLÈME: Méthode '%s' utilise %s receiver, mais '%s' utilise %s (type '%s').

POURQUOI: Cohérence requise:
  - Si une méthode a un pointer receiver, toutes devraient l'avoir
  - Évite les copies accidentelles
  - Simplifie le raisonnement sur les mutations

EXCEPTION: Types immuables (time.Time), maps, funcs, chans.

EXEMPLE INCORRECT:
  func (u User) Name() string     // value
  func (u *User) SetName(n string) // pointer
  func (u User) Validate() error  // value - incohérent!

EXEMPLE CORRECT:
  func (u *User) Name() string
  func (u *User) SetName(n string)
  func (u *User) Validate() error`,
	})

	Register(Message{
		Code:  "KTN-STRUCT-009",
		Short: "receiver '%s' devrait être '%s' (cohérence sur type '%s')",
		Verbose: `PROBLÈME: Nom de receiver '%s' devrait être '%s' (type '%s').

CONVENTION GO:
  - 1-2 lettres, abréviation du type
  - Cohérent sur toutes les méthodes
  - Éviter: me, this, self

Go Code Review Comments:
"The name of a method's receiver should be a reflection
of its identity; often a one or two letter abbreviation
of its type suffices."

EXEMPLE INCORRECT:
  func (u User) Method1() {}
  func (user User) Method2() {}  // Différent!
  func (this User) Method3() {}  // Générique!

EXEMPLE CORRECT:
  func (u User) Method1() {}
  func (u User) Method2() {}
  func (u User) Method3() {}`,
	})

}
