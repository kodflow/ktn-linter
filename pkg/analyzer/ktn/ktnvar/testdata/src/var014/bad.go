// Bad examples for the var015 test case.
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
	// BAD_MAX_ATTEMPTS is placed after var
	BAD_MAX_ATTEMPTS int = 5
	// BAD_DEFAULT_PORT is also after var
	BAD_DEFAULT_PORT int = 8080
)
