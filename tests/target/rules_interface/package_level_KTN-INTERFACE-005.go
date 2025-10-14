// Tests for KTN-INTERFACE-005 (version corrigée)
package rules_interface

// KTN-INTERFACE-005 GOOD: Les interfaces sont maintenant dans interfaces.go

// paymentProcessorImplI005Good est l'implémentation privée.
type paymentProcessorImplI005Good struct {
	gateway string
}

// ProcessPaymentI005Good implémente l'interface.
//
// Params:
//   - amount: le montant à traiter
//
// Returns:
//   - error: une erreur si l'opération échoue
func (p *paymentProcessorImplI005Good) ProcessPaymentI005Good(amount float64) error {
	return nil
}

// NewPaymentProcessorI005Good crée une nouvelle instance.
//
// Params:
//   - gateway: la passerelle de paiement
//
// Returns:
//   - PaymentProcessorI005Good: une nouvelle instance
func NewPaymentProcessorI005Good(gateway string) PaymentProcessorI005Good {
	return &paymentProcessorImplI005Good{gateway: gateway}
}

// emailSenderImplI005Good est l'implémentation privée.
type emailSenderImplI005Good struct {
	smtpServer string
}

// SendEmailI005Good implémente l'interface.
//
// Params:
//   - to: le destinataire
//   - subject: le sujet
//   - body: le corps du message
//
// Returns:
//   - error: une erreur si l'opération échoue
func (e *emailSenderImplI005Good) SendEmailI005Good(to string, subject string, body string) error {
	return nil
}

// NewEmailSenderI005Good crée une nouvelle instance.
//
// Params:
//   - smtpServer: le serveur SMTP
//
// Returns:
//   - EmailSenderI005Good: une nouvelle instance
func NewEmailSenderI005Good(smtpServer string) EmailSenderI005Good {
	return &emailSenderImplI005Good{smtpServer: smtpServer}
}
