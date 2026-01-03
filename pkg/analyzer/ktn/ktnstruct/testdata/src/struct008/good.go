// Package struct008 contains test cases for KTN-STRUCT-008.
package struct008

// User has consistent pointer receivers.
type User struct {
	name  string
	email string
}

// Name uses pointer receiver.
func (u *User) Name() string {
	return u.name
}

// SetName uses pointer receiver.
func (u *User) SetName(n string) {
	u.name = n
}

// Email uses pointer receiver.
func (u *User) Email() string {
	return u.email
}

// SetEmail uses pointer receiver.
func (u *User) SetEmail(e string) {
	u.email = e
}

// Service has consistent value receivers.
type Service struct {
	id   int
	name string
}

// ID uses value receiver.
func (s Service) ID() int {
	return s.id
}

// Name uses value receiver.
func (s Service) Name() string {
	return s.name
}

// SingleMethod has only one method - no consistency check needed.
type SingleMethod struct {
	value int
}

// Value is the only method.
func (s *SingleMethod) Value() int {
	return s.value
}
