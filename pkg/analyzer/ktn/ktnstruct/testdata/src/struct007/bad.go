// Package struct007 contains test cases for KTN-STRUCT-007.
package struct007

// BadUserDTO is a DTO struct with missing serialization tags.
type BadUserDTO struct {
	ID    int    `json:"id"`
	Name  string // want "KTN-STRUCT-007"
	Email string // want "KTN-STRUCT-007"
	age   int    // Private field, no error expected
}

// BadRequestDTO is a DTO with partial tags.
type BadRequestDTO struct {
	ID     int    `json:"id"`
	Method string // want "KTN-STRUCT-007"
	Path   string `xml:"path"`
	Body   string // want "KTN-STRUCT-007"
}

// BadResponseDTO is a DTO with partial tags.
type BadResponseDTO struct {
	Code    int    `json:"code"`
	Status  int    // want "KTN-STRUCT-007"
	Message string // want "KTN-STRUCT-007"
}
