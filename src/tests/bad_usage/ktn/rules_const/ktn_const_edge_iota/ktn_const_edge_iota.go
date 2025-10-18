package badconstiota

// Violations multiples avec iota

const (
	// firstValue sans commentaire de type.
	firstValue int = iota
	// secondValue is the second value.
	secondValue
	// thirdValue is the third value.
	thirdValue
)

// httpStatus mauvais nommage pour HTTP (devrait être HTTPStatus).
const httpStatus int = iota

const (
	// StatusOK commentaire présent.
	StatusOK int = 200 + iota
	// statusBad is a bad status.
	statusBad // Pas de commentaire
	// StatusError is an error status.
	StatusError
)

// Groupe sans commentaire explicatif
const (
	// _ is a blank identifier for iota.
	_ int = iota
	// KB represents one kilobyte.
	KB float64 = 1 << (10 * iota)
	// MB represents one megabyte.
	MB
	// GB represents one gigabyte.
	GB
	// tb represents one terabyte.
	tb // Mauvais nommage (minuscule)
)
