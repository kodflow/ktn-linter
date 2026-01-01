// Package messages provides structured error messages for KTN rules.
// Each rule has a short message (default) and a verbose message (--verbose).
package messages

import (
	"fmt"
	"strings"
)

// initialRegistryCapacity capacité initiale du registre de messages.
const initialRegistryCapacity int = 100

// Message contient les messages court et verbose pour une règle.
// Chaque règle KTN a un message short (affiché par défaut) et un message verbose (avec --verbose).
type Message struct {
	Code    string
	Short   string
	Verbose string
}

// Format retourne toujours le message verbose (détaillé).
// Le paramètre verbose est conservé pour compatibilité mais ignoré.
//
// Params:
//   - verbose: ignoré (conservé pour compatibilité API)
//   - args: arguments pour le formatage
//
// Returns:
//   - string: message formaté (toujours verbose)
func (m Message) Format(_ bool, args ...any) string {
	// Toujours utiliser le message verbose (long) par défaut
	template := m.Verbose
	// Fallback sur Short si pas de message verbose défini
	if template == "" {
		template = m.Short
	}

	// Formater avec les arguments
	if len(args) > 0 {
		// Appliquer le formatage
		return fmt.Sprintf(template, args...)
	}

	// Retour du template sans formatage
	return template
}

// FormatShort retourne le message court avec suffixe --verbose.
//
// Params:
//   - args: arguments pour le formatage
//
// Returns:
//   - string: message court formaté
func (m Message) FormatShort(args ...any) string {
	msg := m.Format(false, args...)
	// Ajouter le suffixe si verbose disponible
	if m.Verbose != "" && !strings.HasSuffix(msg, "(--verbose pour détails)") {
		msg += " (--verbose pour détails)"
	}
	// Retour du message
	return msg
}

// FormatVerbose retourne le message verbose complet.
//
// Params:
//   - args: arguments pour le formatage
//
// Returns:
//   - string: message verbose formaté
func (m Message) FormatVerbose(args ...any) string {
	// Retour du message verbose
	return m.Format(true, args...)
}

// registry stocke tous les messages par code de règle.
var registry map[string]Message = make(map[string]Message, initialRegistryCapacity)

// Register enregistre un message pour une règle.
//
// Params:
//   - msg: message à enregistrer
func Register(msg Message) {
	registry[msg.Code] = msg
}

// Get récupère un message par code de règle.
//
// Params:
//   - code: code de la règle (ex: KTN-FUNC-001)
//
// Returns:
//   - Message: message trouvé ou vide
//   - bool: true si trouvé
func Get(code string) (Message, bool) {
	msg, ok := registry[code]
	// Retour du résultat
	return msg, ok
}

// init enregistre tous les messages.
func init() {
	// Enregistrer tous les messages
	registerAPIMessages()
	registerCommentMessages()
	registerConstMessages()
	registerFuncMessages()
	registerGenericMessages()
	registerStructMessages()
	registerTestMessages()
	registerVarMessages()
	registerInterfaceMessages()
	registerReturnMessages()
}
