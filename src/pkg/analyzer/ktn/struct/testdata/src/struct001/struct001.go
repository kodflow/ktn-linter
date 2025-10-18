package struct001

// CORRECT: Structs en MixedCaps

// UserConfig utilise MixedCaps.
type UserConfig struct {
	Name string
}

// httpClient utilise mixedCaps (privé).
type httpClient struct {
	URL string
}

// HTTPServer utilise un initialisme.
type HTTPServer struct {
	Port int
}

// BAD: Structs avec underscores

// user_config utilise snake_case (incorrect).
type user_config struct { // want "KTN-STRUCT-001.*MixedCaps"
	name string
}

// User_Profile mélange majuscules et underscores (incorrect).
type User_Profile struct { // want "KTN-STRUCT-001.*MixedCaps"
	ID int
}

// api_response utilise snake_case (incorrect).
type api_response struct { // want "KTN-STRUCT-001.*MixedCaps"
	data string
}
