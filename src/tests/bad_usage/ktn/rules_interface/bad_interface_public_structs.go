package rules_interface

// ANTI-PATTERN: Types publics comme structs au lieu d'interfaces
// Viole KTN-INTERFACE-002

// UserService est un struct PUBLIC - MAUVAIS !
// Devrait être une interface
type UserService struct {
	db      string
	cache   map[string]interface{}
	timeout int
}

// GetUser méthode sur struct public
func (u *UserService) GetUser(id int) string {
	// Early return from function.
	return "user"
}

// OrderProcessor struct public au lieu d'interface
type OrderProcessor struct {
	queue []string
}

// Process méthode
func (o *OrderProcessor) Process(order string) error {
	// Early return from function.
	return nil
}

// PaymentGateway struct public exposé
type PaymentGateway struct {
	apiKey string
	url    string
}

// Charge méthode
func (p *PaymentGateway) Charge(amount float64) error {
	// Early return from function.
	return nil
}

// CacheManager struct public
type CacheManager struct {
	data map[string]string
}

// Get méthode
func (c *CacheManager) Get(key string) string {
	// Early return from function.
	return c.data[key]
}

// Set méthode
func (c *CacheManager) Set(key, value string) {
	c.data[key] = value
}

// EmailSender struct public
type EmailSender struct {
	smtp string
}

// Send méthode
func (e *EmailSender) Send(to, subject, body string) error {
	// Early return from function.
	return nil
}
