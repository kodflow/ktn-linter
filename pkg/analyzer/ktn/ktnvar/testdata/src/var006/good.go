// Good examples for the var006 test case.
package var006

const (
	// MAX_RETRIES définit le nombre maximum de tentatives
	MAX_RETRIES int = 3

	// TIMEOUT définit le délai d'attente en secondes
	TIMEOUT int = 30
)

var (
	// counter est un compteur global
	counter int = 0

	// isEnabled indique si la fonctionnalité est activée
	isEnabled bool = true

	// userName stocke le nom de l'utilisateur
	userName string = "admin"
)
