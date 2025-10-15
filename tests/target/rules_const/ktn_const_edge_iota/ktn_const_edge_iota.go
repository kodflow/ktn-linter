package rules_const

// Utilisation correcte de iota avec documentation complète

const (
	// FirstValue première valeur de l'énumération.
	FirstValue = iota
	// SecondValue deuxième valeur de l'énumération.
	SecondValue
	// ThirdValue troisième valeur de l'énumération.
	ThirdValue
)

// HTTPStatus représente un code de statut HTTP.
const HTTPStatus = iota

const (
	// StatusOK code HTTP 200 - succès.
	StatusOK = 200 + iota
	// StatusBad code HTTP 400 - requête invalide.
	StatusBad
	// StatusError code HTTP 500 - erreur serveur.
	StatusError
)

// Unités de stockage en bytes.
const (
	_  = iota // Ignore la première valeur
	// KB représente un kilobyte (1024 bytes).
	KB float64 = 1 << (10 * iota)
	// MB représente un megabyte.
	MB
	// GB représente un gigabyte.
	GB
	// TB représente un terabyte.
	TB
)
