package badconstiota

// Violations multiples avec iota

const (
	// firstValue sans commentaire de type
	firstValue = iota
	secondValue
	thirdValue
)

// httpStatus mauvais nommage pour HTTP (devrait être HTTPStatus)
const httpStatus = iota

const (
	// StatusOK commentaire présent
	StatusOK = 200 + iota
	statusBad // Pas de commentaire
	StatusError
)

// Groupe sans commentaire explicatif
const (
	_          = iota
	KB float64 = 1 << (10 * iota)
	MB
	GB
	tb // Mauvais nommage (minuscule)
)
