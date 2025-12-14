// Const messages for KTN-CONST rules.
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
}
