package struct003

// CORRECT: Champs exportés documentés

// GoodConfig représente une configuration.
type GoodConfig struct {
	// Name est le nom de la configuration.
	Name string

	// Port est le port d'écoute.
	Port int

	// Les champs privés n'ont pas besoin de documentation
	internal string
}

// BAD: Champs exportés non documentés

// BadConfig représente une configuration.
type BadConfig struct {
	Name string // want "KTN-STRUCT-003.*sans commentaire"

	Port int // want "KTN-STRUCT-003.*sans commentaire"

	// Champ privé (pas d'erreur)
	internal string
}

// AnotherBadStruct avec champs non documentés.
type AnotherBadStruct struct {
	ID   int    // want "KTN-STRUCT-003.*sans commentaire"
	Data string // want "KTN-STRUCT-003.*sans commentaire"
}
