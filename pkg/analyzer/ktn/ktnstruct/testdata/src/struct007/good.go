// Package struct007 contains test cases for KTN-STRUCT-007.
package struct007

// UserDTO is a properly tagged DTO struct.
type UserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	age   int    // Private field, no tag needed
}

// RequestDTO is a properly tagged DTO.
type RequestDTO struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Path   string `xml:"path"`
	Body   string `json:"body"`
}

// ResponseDTO is a properly tagged DTO.
type ResponseDTO struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

// NotADTO is a regular struct (no DTO suffix, no serialization tags).
type User struct {
	Name  string
	Email string
}

// Settings is a regular struct without DTO naming.
type Settings struct {
	host string
	port int
}
