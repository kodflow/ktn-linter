package struct002

// CORRECT: Structs avec documentation godoc

// UserConfig représente la configuration utilisateur.
type UserConfig struct {
	Name string
}

// HTTPClient représente un client HTTP.
type HTTPClient struct {
	URL string
}

// BAD: Structs sans documentation godoc

type BadNoDoc struct { // want "KTN-STRUCT-002.*commentaire godoc"
	Field string
}

type AnotherBadStruct struct { // want "KTN-STRUCT-002.*commentaire godoc"
	ID   int
	Name string
}
