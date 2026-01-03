// Package struct008 contains test cases for KTN-STRUCT-008.
package struct008

// BadUser has inconsistent receiver types.
type BadUser struct {
	name  string
	email string
}

// Name uses pointer receiver.
func (u *BadUser) Name() string {
	return u.name
}

// SetName uses pointer receiver.
func (u *BadUser) SetName(n string) {
	u.name = n
}

// Email uses value receiver - inconsistent! // want "KTN-STRUCT-008"
func (u BadUser) Email() string {
	return u.email
}

// BadService has mixed receiver types.
type BadService struct {
	id   int
	name string
}

// ID uses value receiver.
func (s BadService) ID() int {
	return s.id
}

// Name uses value receiver.
func (s BadService) Name() string {
	return s.name
}

// SetID uses pointer receiver - inconsistent! // want "KTN-STRUCT-008"
func (s *BadService) SetID(id int) {
	s.id = id
}

// SetName uses pointer receiver - inconsistent! // want "KTN-STRUCT-008"
func (s *BadService) SetName(name string) {
	s.name = name
}
