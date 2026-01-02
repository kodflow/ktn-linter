// Package messages provides structured error messages for KTN rules.
// This file contains CONST rule messages.
package messages

// registerConstMessages enregistre les messages CONST.
func registerConstMessages() {
	Register(Message{
		Code:  "KTN-CONST-001",
		Short: "constante '%s' sans type explicite",
		Verbose: `PROBLÈME: La constante '%s' n'a pas de type déclaré.

POURQUOI: Le typage explicite:
  - Documente l'intention
  - Évite les conversions implicites
  - Améliore la lisibilité

EXEMPLE INCORRECT:
  const MaxSize = 1024
  const Prefix = "app_"

EXEMPLE CORRECT:
  const MaxSize int = 1024
  const Prefix string = "app_"

EXCEPTION: Les constantes iota peuvent omettre le type.`,
	})

	Register(Message{
		Code:  "KTN-CONST-002",
		Short: "constantes mal organisées. Ordre: const → var → type → func",
		Verbose: `PROBLÈME: Les constantes ne sont pas groupées ou mal placées.

POURQUOI: L'ordre standardisé facilite la navigation:
  1. const (constantes)
  2. var (variables de package)
  3. type (types et structs)
  4. func (fonctions)

EXEMPLE INCORRECT:
  func DoSomething() {}
  const MaxSize = 1024  // Après une fonction!

EXEMPLE CORRECT:
  const (
      MaxSize  int    = 1024
      MinSize  int    = 64
  )

  var defaultConfig = Config{}

  type Config struct { ... }

  func DoSomething() {}`,
	})

	Register(Message{
		Code:  "KTN-CONST-003",
		Short: "constante '%s' mal nommée. CamelCase requis, pas SCREAMING_SNAKE",
		Verbose: `PROBLÈME: La constante '%s' utilise SCREAMING_SNAKE_CASE.

POURQUOI: Go utilise CamelCase pour TOUT, y compris constantes.
SCREAMING_SNAKE_CASE est une convention C/Java, pas Go.

EXEMPLE INCORRECT:
  const MAX_BUFFER_SIZE = 1024
  const DEFAULT_TIMEOUT = 30

EXEMPLE CORRECT:
  const MaxBufferSize = 1024
  const DefaultTimeout = 30

NOTE: Constantes privées = camelCase (minuscule initiale).
  const maxBufferSize = 1024  // Privée
  const MaxBufferSize = 1024  // Exportée`,
	})

	Register(Message{
		Code:  "KTN-CONST-004",
		Short: "constante '%s' trop courte (min %d caractères)",
		Verbose: `PROBLÈME: La constante '%s' a un nom trop court.

POURQUOI: Les noms courts manquent de contexte et de clarté.
Un minimum de %d caractères assure une meilleure lisibilité.

EXCEPTION: Le blank identifier (_) est toujours autorisé.

EXEMPLE INCORRECT:
  const A int = 1
  const B int = 2

EXEMPLE CORRECT:
  const MaxRetries int = 1
  const MinSize int = 2
  const ID int = 1  // 2 caractères OK`,
	})

	Register(Message{
		Code:  "KTN-CONST-005",
		Short: "constante '%s' trop longue (%d chars, max %d)",
		Verbose: `PROBLÈME: La constante '%s' a un nom trop long (%d caractères).

POURQUOI: Les noms trop longs réduisent la lisibilité.
Maximum recommandé: %d caractères.

EXEMPLE INCORRECT:
  const DefaultHTTPConnectionTimeoutInSeconds = 30

EXEMPLE CORRECT:
  const HTTPTimeout = 30
  const DefaultConnTimeout = 30`,
	})

	Register(Message{
		Code:  "KTN-CONST-006",
		Short: "constante '%s' masque un identifiant built-in",
		Verbose: `PROBLÈME: La constante '%s' masque un identifiant built-in Go.

POURQUOI: Masquer les built-ins cause des comportements inattendus:
  - Confusion pour les lecteurs du code
  - Erreurs subtiles difficiles à déboguer
  - Le built-in devient inaccessible dans ce scope

BUILT-INS PROTÉGÉS:
  Types: bool, byte, int, string, error, any, ...
  Constantes: true, false, iota
  Fonctions: len, cap, make, new, append, panic, ...
  Zero-value: nil

EXEMPLE INCORRECT:
  const len int = 100   // Masque len()
  const nil int = 0     // Masque nil

EXEMPLE CORRECT:
  const MaxLen int = 100
  const NilValue int = 0`,
	})
}
