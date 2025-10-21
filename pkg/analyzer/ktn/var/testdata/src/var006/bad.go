package var006

// Vars come first (good placement)
var (
	// badCounter est un compteur global
	badCounter int = 0
	// badStatus stocke le statut actuel
	badStatus string = "idle"
)

// First const block AFTER vars (violates VAR-006 #1)
const (
	// BAD_MAX_ATTEMPTS is placed after var
	BAD_MAX_ATTEMPTS int = 5 // want "KTN-VAR-006: les constantes doivent être déclarées avant les variables"
)

// Second const block AFTER vars (violates VAR-006 #2)
const (
	// BAD_DEFAULT_PORT is also after var
	BAD_DEFAULT_PORT int = 8080 // want "KTN-VAR-006: les constantes doivent être déclarées avant les variables"
)
