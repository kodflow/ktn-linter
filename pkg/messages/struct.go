// Package messages provides structured error messages for KTN rules.
// This file contains STRUCT rule messages.
package messages

// registerStructMessages enregistre les messages STRUCT.
func registerStructMessages() {
	Register(Message{
		Code:  "KTN-STRUCT-001",
		Short: "struct '%s' a %d méthode(s) publique(s) sans interface",
		Verbose: `PROBLÈME: La struct '%s' a %d méthode(s) sans interface.

POURQUOI: Une interface permet:
  - Le mocking dans les tests
  - L'injection de dépendances
  - Le découplage

SOLUTIONS:
  1. Interface locale: créer dans le même fichier
  2. Cross-package: var _ port.I = (*%s)(nil)

EXCEPTIONS AUTOMATIQUES:
  - DTOs (tags json/xml/yaml)
  - Consumers (champs de type interface)

EXEMPLE SOLUTION 1:
  type UserRepository interface {
      GetByID(id int) (*User, error)
  }
  type userRepositoryImpl struct { ... }

EXEMPLE SOLUTION 2 (cross-package):
  var _ port.UserRepository = (*UserRepositoryImpl)(nil)`,
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
}
