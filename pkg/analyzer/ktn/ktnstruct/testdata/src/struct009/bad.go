// Package struct009 contains test cases for KTN-STRUCT-009.
package struct009

// BadUser is a struct with inconsistent receiver names.
type BadUser struct {
	name string
}

// Name uses receiver 'u'.
func (u *BadUser) Name() string {
	return u.name
}

// SetName uses receiver 'user' - inconsistent! // want "KTN-STRUCT-009"
func (user *BadUser) SetName(n string) {
	user.name = n
}

// BadService uses generic receiver names.
type BadService struct {
	id int
}

// ID uses 'this' - generic name! // want "KTN-STRUCT-009"
func (this *BadService) ID() int {
	return this.id
}

// SetID uses 'self' - generic name! // want "KTN-STRUCT-009"
func (self *BadService) SetID(id int) {
	self.id = id
}

// BadClient has mixed receiver names.
type BadClient struct {
	connected bool
}

// IsConnected uses 'c'.
func (c *BadClient) IsConnected() bool {
	return c.connected
}

// Connect uses 'cl' - inconsistent! // want "KTN-STRUCT-009"
func (cl *BadClient) Connect() {
	cl.connected = true
}

// Disconnect uses 'client' - inconsistent! // want "KTN-STRUCT-009"
func (client *BadClient) Disconnect() {
	client.connected = false
}
