// Package gospec_good_naming démontre les conventions de nommage idiomatiques de Go.
// Référence: https://go.dev/doc/effective_go#names
// Référence: https://github.com/golang/go/wiki/CodeReviewComments#naming
package gospec_good_naming

import "fmt"

// ✅ GOOD: Using MixedCaps for exported names
var UserCount int
var MaxRetryCount int

// ✅ GOOD: Using mixedCaps for unexported names
var defaultTimeout = 30
var apiVersion = "v1"

// ✅ GOOD: Function names using MixedCaps
func GetUserName() string { return "" }
func CalculateTotalPrice() int { return 0 }

// ✅ GOOD: Type names using MixedCaps
type UserData struct {
	FirstName string
	LastName  string
}

// ✅ GOOD: Grouping related constants
const (
	StatusPending  = "pending"
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// ✅ GOOD: Clear, descriptive names
func ValidateUserInput(input string) error {
	if input == "" {
		return fmt.Errorf("empty input")
	}
	return nil
}

func FetchUserByID(id int) (*UserData, error) {
	return &UserData{}, nil
}

// ✅ GOOD: Short names in limited scope
func ProcessItems(items []int) int {
	sum := 0
	for _, v := range items { // 'v' acceptable dans scope limité
		sum += v
	}
	return sum
}

// ✅ GOOD: Conventional abbreviations
var cfg Config
var id int
var db Database

// ✅ GOOD: Consistent receiver names
type Handler struct {
	count int
}

func (h Handler) Handle() {}
func (h Handler) Execute() {}
func (h Handler) Run() {}

// ✅ GOOD: Interface with 'er' suffix for single method
type Reader interface {
	Read() ([]byte, error)
}

type Writer interface {
	Write([]byte) error
}

// ✅ GOOD: Boolean names
var Enabled bool
var IsValid bool
var HasPermission bool

// ✅ GOOD: Proper acronym casing
var HTTPClient int
var XMLParser int
var JSONData int
var IDValue int

// ✅ GOOD: Getter without 'Get' prefix (when simple accessor)
func (h Handler) Count() int {
	return h.count
}

// ✅ GOOD: Setter with 'Set' prefix
func (h *Handler) SetCount(n int) {
	h.count = n
}

// ✅ GOOD: Error variables with 'Err' prefix
var ErrNotFound = fmt.Errorf("not found")
var ErrInvalidInput = fmt.Errorf("invalid input")
var ErrTimeout = fmt.Errorf("timeout")

// ✅ GOOD: Package-level types with clear purpose
type Config struct {
	Host    string
	Port    int
	Timeout int
}

type Database interface {
	Query(string) ([]byte, error)
	Close() error
}

// ✅ GOOD: Constructor function
func NewHandler() *Handler {
	return &Handler{count: 0}
}

func NewConfig(host string, port int) *Config {
	return &Config{
		Host: host,
		Port: port,
	}
}

// ✅ GOOD: Method names are verbs
func (h *Handler) Increment() {
	h.count++
}

func (h *Handler) Reset() {
	h.count = 0
}

// ✅ GOOD: Clear parameter names
func ProcessUser(name string, age int, email string) error {
	_, _, _ = name, age, email
	return nil
}

// ✅ GOOD: Using conventional names for common patterns
func (h Handler) String() string { // Stringer interface
	return fmt.Sprintf("Handler{count: %d}", h.count)
}

// ✅ GOOD: Unexported helper functions
func validateInput(s string) bool {
	return s != ""
}

func parseData(data []byte) (string, error) {
	return string(data), nil
}

// ✅ GOOD: Package name matches directory (assumed: naming/)
// package naming would be in directory 'naming'
