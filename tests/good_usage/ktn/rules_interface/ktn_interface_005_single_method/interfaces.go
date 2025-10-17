// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-005: Interfaces dans interfaces.go (✅ CORRIGÉ)
// ════════════════════════════════════════════════════════════════════════════

package KTN_INTERFACE_005

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
