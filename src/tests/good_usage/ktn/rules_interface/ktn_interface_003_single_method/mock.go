//go:build test
// +build test

package KTN_INTERFACE_005

// MockPaymentProcessorI005Good est le mock de PaymentProcessorI005Good.
type MockPaymentProcessorI005Good struct {
	ProcessPaymentI005GoodFunc func(amount float64) error
}

// ProcessPaymentI005Good implémente l'interface PaymentProcessorI005Good.
//
// Params:
//   - amount: le montant à traiter
//
// Returns:
//   - error: une erreur si l'opération échoue
func (m *MockPaymentProcessorI005Good) ProcessPaymentI005Good(amount float64) error {
	if m.ProcessPaymentI005GoodFunc != nil {
		// Early return from function.
		return m.ProcessPaymentI005GoodFunc(amount)
	}
	// Early return from function.
	return nil
}

// MockEmailSenderI005Good est le mock de EmailSenderI005Good.
type MockEmailSenderI005Good struct {
	SendEmailI005GoodFunc func(to string, subject string, body string) error
}

// SendEmailI005Good implémente l'interface EmailSenderI005Good.
//
// Params:
//   - to: le destinataire
//   - subject: le sujet
//   - body: le corps du message
//
// Returns:
//   - error: une erreur si l'opération échoue
func (m *MockEmailSenderI005Good) SendEmailI005Good(to string, subject string, body string) error {
	if m.SendEmailI005GoodFunc != nil {
		// Early return from function.
		return m.SendEmailI005GoodFunc(to, subject, body)
	}
	// Early return from function.
	return nil
}
