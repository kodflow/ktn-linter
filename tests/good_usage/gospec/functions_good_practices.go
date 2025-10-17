// Package gospec_good_functions démontre les patterns de fonctions idiomatiques.
// Référence: https://go.dev/doc/effective_go#functions
package gospec_good_functions

import (
	"fmt"
	"time"
)

// ✅ GOOD: Clear, focused function with single responsibility
func ValidateInput(input string) error {
	if input == "" {
		return fmt.Errorf("empty input")
	}
	return nil
}

// ✅ GOOD: Using struct for multiple parameters
type UserConfig struct {
	Name    string
	Email   string
	Age     int
	Address string
}

func CreateUser(cfg UserConfig) (*User, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("name required")
	}
	return &User{cfg.Name, cfg.Email}, nil
}

// ✅ GOOD: Functional options pattern
type ServerOption func(*Server)

func WithTimeout(d time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = d
	}
}

func WithHost(host string) ServerOption {
	return func(s *Server) {
		s.host = host
	}
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		host:    "localhost",
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ✅ GOOD: Consistent parameter order across related functions
func ReadFile(filename string) ([]byte, error) {
	return nil, nil
}

func WriteFile(filename string, data []byte) error {
	_, _ = filename, data
	return nil
}

// ✅ GOOD: Using variadic parameters appropriately
func Sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// ✅ GOOD: Clear return values
func GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id")
	}
	return &User{Name: "test", Email: "test@example.com"}, nil
}

// ✅ GOOD: Named returns for clarity in complex functions
func Calculate(a, b int) (sum, product int, err error) {
	if a < 0 || b < 0 {
		err = fmt.Errorf("negative values not allowed")
		return
	}
	sum = a + b
	product = a * b
	return
}

// ✅ GOOD: Consistent receiver type (all pointer or all value)
type Counter struct {
	count int
}

func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) Decrement() {
	c.count--
}

func (c *Counter) Value() int {
	return c.count
}

// ✅ GOOD: Constructor with validation
func NewCounter(initial int) (*Counter, error) {
	if initial < 0 {
		return nil, fmt.Errorf("initial count cannot be negative")
	}
	return &Counter{count: initial}, nil
}

// ✅ GOOD: Simple constructor for zero value initialization
func NewDefaultCounter() *Counter {
	return &Counter{}
}

// ✅ GOOD: Method clearly uses receiver
func (c *Counter) Reset() {
	c.count = 0
}

// ✅ GOOD: Exported function with clear documentation need
// ProcessData processes the input data and returns results.
func ProcessData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	return data, nil
}

// ✅ GOOD: Function composition
func Validate(validators ...func(string) error) func(string) error {
	return func(s string) error {
		for _, validator := range validators {
			if err := validator(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func NotEmpty(s string) error {
	if s == "" {
		return fmt.Errorf("empty string")
	}
	return nil
}

func MinLength(min int) func(string) error {
	return func(s string) error {
		if len(s) < min {
			return fmt.Errorf("too short")
		}
		return nil
	}
}

// ✅ GOOD: Defer for cleanup
func ProcessFile(filename string) error {
	f, err := openResource(filename)
	if err != nil {
		return err
	}
	defer closeResource(f)

	return processResource(f)
}

// ✅ GOOD: Higher-order function
func Map(items []int, fn func(int) int) []int {
	result := make([]int, len(items))
	for i, v := range items {
		result[i] = fn(v)
	}
	return result
}

func Filter(items []int, fn func(int) bool) []int {
	result := make([]int, 0)
	for _, v := range items {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// ✅ GOOD: Closure with clear purpose
func Counter2() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// ✅ GOOD: Table-driven approach via function
func TestValidator() {
	tests := []struct {
		input string
		valid bool
	}{
		{"valid", true},
		{"", false},
	}

	for _, tt := range tests {
		err := ValidateInput(tt.input)
		if (err == nil) != tt.valid {
			fmt.Printf("unexpected result for %q\n", tt.input)
		}
	}
}

// ✅ GOOD: Init function is minimal
func init() {
	// Only essential initialization
	setupDefaults()
}

// ✅ GOOD: Separating concerns into multiple functions
func SaveUser(u *User) error {
	if err := validateUser(u); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := persistUser(u); err != nil {
		return fmt.Errorf("persistence failed: %w", err)
	}

	notifyUserCreated(u)
	return nil
}

func validateUser(u *User) error {
	if u.Name == "" {
		return fmt.Errorf("name required")
	}
	return nil
}

func persistUser(u *User) error {
	// Save to database
	return nil
}

func notifyUserCreated(u *User) {
	// Send notification
}

// ✅ GOOD: Using interfaces for flexibility
type Storage interface {
	Save(key string, value []byte) error
	Load(key string) ([]byte, error)
}

func ProcessWithStorage(s Storage, key string) error {
	data, err := s.Load(key)
	if err != nil {
		return err
	}
	return s.Save(key, data)
}

// ✅ GOOD: Small, focused helper functions
func isValidEmail(email string) bool {
	return len(email) > 0 && containsAt(email)
}

func containsAt(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}

// Helper types
type User struct {
	Name  string
	Email string
}

type Server struct {
	host    string
	timeout time.Duration
}

func openResource(string) (int, error)  { return 0, nil }
func closeResource(int)                 {}
func processResource(int) error         { return nil }
func setupDefaults()                    {}
