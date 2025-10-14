// interfaces.go pour les règles INTERFACE (version corrigée)
package rules_interface

// KTN-INTERFACE-001 GOOD: Ce fichier existe, donc la règle est respectée

// UserServiceI001Good définit le contrat du service utilisateur.
type UserServiceI001Good interface {
	// GetUser récupère un utilisateur.
	//
	// Params:
	//   - id: l'identifiant de l'utilisateur
	//
	// Returns:
	//   - string: le nom de l'utilisateur
	//   - error: une erreur si l'opération échoue
	GetUser(id string) (string, error)
}

// CacheManagerI006Good définit le contrat de gestion du cache.
type CacheManagerI006Good interface {
	// GetI006Good récupère une valeur du cache.
	//
	// Params:
	//   - key: la clé à rechercher
	//
	// Returns:
	//   - string: la valeur trouvée
	//   - bool: true si la clé existe
	GetI006Good(key string) (string, bool)

	// SetI006Good définit une valeur dans le cache.
	//
	// Params:
	//   - key: la clé
	//   - value: la valeur
	SetI006Good(key string, value string)
}

// LoggerI006Good définit le contrat de journalisation.
type LoggerI006Good interface {
	// InfoI006Good enregistre un message d'information.
	//
	// Params:
	//   - msg: le message à enregistrer
	InfoI006Good(msg string)

	// ErrorI006Good enregistre un message d'erreur.
	//
	// Params:
	//   - msg: le message d'erreur
	ErrorI006Good(msg string)
}

// PaymentProcessorI005Good définit le contrat de traitement des paiements.
type PaymentProcessorI005Good interface {
	// ProcessPaymentI005Good traite un paiement.
	//
	// Params:
	//   - amount: le montant à traiter
	//
	// Returns:
	//   - error: une erreur si l'opération échoue
	ProcessPaymentI005Good(amount float64) error
}

// EmailSenderI005Good définit le contrat d'envoi d'emails.
type EmailSenderI005Good interface {
	// SendEmailI005Good envoie un email.
	//
	// Params:
	//   - to: le destinataire
	//   - subject: le sujet
	//   - body: le corps du message
	//
	// Returns:
	//   - error: une erreur si l'opération échoue
	SendEmailI005Good(to string, subject string, body string) error
}

// MarkerInterfaceI006Good est une interface marqueur.
type MarkerInterfaceI006Good interface{}
