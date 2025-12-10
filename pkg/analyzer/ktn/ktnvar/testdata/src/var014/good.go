// Good examples for the var015 test case.
package var014

const (
	// MaxRetries définit le nombre maximum de tentatives
	MaxRetries int = 3

	// Timeout définit le délai d'attente en secondes
	Timeout int = 30
)

var (
	// counter est un compteur global
	counter int = 0

	// isEnabled indique si la fonctionnalité est activée
	isEnabled bool = true

	// userName stocke le nom de l'utilisateur
	userName string = "admin"
)
