// Package var014 contains test cases for KTN rules.
package var014

// Vars come first (good placement)
var (
	// badCounter est un compteur global
	badCounter int = 0
	// badStatus stocke le statut actuel
	badStatus string = "idle"
)

// Const block AFTER vars (violates VAR-006 and CONST-002 - inévitable)
const (
	// BadMaxAttempts is placed after var
	BadMaxAttempts int = 5
	// BadDefaultPort is also after var
	BadDefaultPort int = 8080
)
