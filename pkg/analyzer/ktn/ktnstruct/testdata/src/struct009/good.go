// Package struct009 contains test cases for KTN-STRUCT-009.
package struct009

// User has consistent receiver names.
type User struct {
	name  string
	email string
}

// Name returns the user's name.
func (u *User) Name() string {
	return u.name
}

// SetName sets the user's name.
func (u *User) SetName(n string) {
	u.name = n
}

// Email returns the user's email.
func (u *User) Email() string {
	return u.email
}

// SetEmail sets the user's email.
func (u *User) SetEmail(e string) {
	u.email = e
}

// Service has consistent receiver names.
type Service struct {
	id   int
	name string
}

// ID returns the service ID.
func (s *Service) ID() int {
	return s.id
}

// Name returns the service name.
func (s *Service) Name() string {
	return s.name
}

// Client has consistent receiver names.
type Client struct {
	connected bool
	timeout   int
}

// IsConnected returns connection status.
func (c *Client) IsConnected() bool {
	return c.connected
}

// Connect establishes connection.
func (c *Client) Connect() {
	c.connected = true
}

// Disconnect closes connection.
func (c *Client) Disconnect() {
	c.connected = false
}

// Timeout returns the timeout value.
func (c *Client) Timeout() int {
	return c.timeout
}

// SingleMethod has only one method - no consistency check needed.
type SingleMethod struct {
	value int
}

// Value returns the value.
func (s *SingleMethod) Value() int {
	return s.value
}
