// ════════════════════════════════════════════════════════════════════════════
// KTN-INTERFACE-003: Interface définie hors de interfaces.go
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//
//	Les interfaces publiques doivent être définies dans interfaces.go
//	pour centraliser tous les contrats du package en un seul endroit.
//
//	POURQUOI :
//	- Facilite la découverte des contrats disponibles
//	- Centralise la documentation des interfaces
//	- Convention standard pour packages Go bien structurés
//	- Évite la dispersion des définitions d'interfaces
//
// ❌ CAS INCORRECT 1: Interface hors interfaces.go (SEULE ERREUR: KTN-INTERFACE-003)
// NOTE: PaymentProcessorI005 devrait être dans interfaces.go
// ERREUR ATTENDUE: KTN-INTERFACE-003 sur PaymentProcessorI005
//
// ❌ CAS INCORRECT 2: Interface hors interfaces.go (SEULE ERREUR: KTN-INTERFACE-003)
// NOTE: EmailSenderI005 devrait être dans interfaces.go
// ERREUR ATTENDUE: KTN-INTERFACE-003 sur EmailSenderI005
//
// ✅ CAS PARFAIT (voir target/) :
//
//	// interfaces.go
//	type PaymentProcessor interface {
//	    ProcessPayment(amount float64) error
//	}
//
//	// impl.go (avec implémentation privée)
//	type paymentProcessorImpl struct { ... }
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_INTERFACE_005

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
	// Early return from function.
	return nil
}
