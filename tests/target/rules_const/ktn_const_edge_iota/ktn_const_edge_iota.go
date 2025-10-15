package rules_const

// Utilisation correcte de iota avec documentation complète.
// Exception KTN-CONST-004: Le type peut être omis avec iota (inféré automatiquement).

// Énumération de valeurs séquentielles avec iota.
const (
	// FirstValue première valeur de l'énumération.
	FirstValue = iota
	// SecondValue deuxième valeur de l'énumération.
	SecondValue
	// ThirdValue troisième valeur de l'énumération.
	ThirdValue
)

// Codes de statut HTTP avec offset.
const (
	// StatusOK code HTTP 200 - succès.
	StatusOK = 200 + iota
	// StatusBad code HTTP 201 - requête invalide.
	StatusBad
	// StatusError code HTTP 202 - erreur serveur.
	StatusError
)

// Unités de stockage en bytes avec type explicite sur première constante.
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
