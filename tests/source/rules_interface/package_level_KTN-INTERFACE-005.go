// Tests for KTN-INTERFACE-005: Interface définie hors de interfaces.go
package rules_interface

// KTN-INTERFACE-005: PaymentProcessorI005 est définie ici au lieu de interfaces.go
// Les interfaces publiques doivent être dans interfaces.go

// PaymentProcessorI005 définit le contrat de traitement des paiements.
// Cette interface devrait être dans interfaces.go.
type PaymentProcessorI005 interface {
	// ProcessPaymentI005 traite un paiement.
	//
	// Params:
	//   - amount: le montant à traiter
	//
	// Returns:
	//   - error: une erreur si l'opération échoue
	ProcessPaymentI005(amount float64) error
}

// EmailSenderI005 définit le contrat d'envoi d'emails.
// Cette interface devrait aussi être dans interfaces.go.
type EmailSenderI005 interface {
	// SendEmailI005 envoie un email.
	//
	// Params:
	//   - to: le destinataire
	//   - subject: le sujet
	//   - body: le corps du message
	//
	// Returns:
	//   - error: une erreur si l'opération échoue
	SendEmailI005(to string, subject string, body string) error
}

// paymentProcessorImplI005 est l'implémentation privée.
type paymentProcessorImplI005 struct {
	gateway string
}

// ProcessPaymentI005 implémente l'interface.
//
// Params:
//   - amount: le montant à traiter
//
// Returns:
//   - error: une erreur si l'opération échoue
func (p *paymentProcessorImplI005) ProcessPaymentI005(amount float64) error {
	return nil
}
