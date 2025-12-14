// Package api001 contains test cases for KTN-API-001.
package api001

import (
	"context"
	"io"
	"net/http"
	"time"
)

// httpDoer is a minimal consumer-side interface.
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// reader is a minimal consumer-side interface for reading.
type reader interface {
	Read(p []byte) (n int, err error)
}

// closer is a minimal consumer-side interface for closing.
type closer interface {
	Close() error
}

// goodWithInterface uses an interface parameter - no warning.
func goodWithInterface(d httpDoer) (*http.Response, error) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	return d.Do(req)
}

// goodWithReaderInterface uses reader interface - no warning.
func goodWithReaderInterface(r reader) ([]byte, error) {
	buf := make([]byte, 100)
	_, err := r.Read(buf)
	// Vérification de la condition
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// goodWithCloserInterface uses closer interface - no warning.
func goodWithCloserInterface(c closer) error {
	return c.Close()
}

// goodWithStdInterface uses standard io.Reader interface - no warning.
func goodWithStdInterface(r io.Reader) ([]byte, error) {
	buf := make([]byte, 100)
	_, err := r.Read(buf)
	// Vérification de la condition
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// goodExternalConcreteNoMethodCalls has external concrete type but no method calls.
// No warning because no methods are called on the parameter.
func goodExternalConcreteNoMethodCalls(client *http.Client) string {
	return "using client with timeout: " + client.Timeout.String()
}

// goodExternalConcreteFieldAccess accesses a field on external type.
// No warning because no methods are called on the parameter directly.
func goodExternalConcreteFieldAccess(req *http.Request) string {
	// Only field access, no method calls on req
	return req.Method + " " + req.URL.String()
}

// goodWithUnusedParam has external concrete type but parameter is unused.
// No warning because no methods are called.
func goodWithUnusedParam(_ *http.Client) string {
	return "client not used"
}

// localType is a type defined in the same package.
type localType struct {
	Value string
}

// LocalMethod is a method on localType.
func (l *localType) LocalMethod() string {
	return l.Value
}

// goodSamePackageType uses a same-package type with method calls - no warning.
func goodSamePackageType(lt *localType) string {
	return lt.LocalMethod()
}

// goodAllowedTimeType uses time.Time which is allowlisted - no warning.
func goodAllowedTimeType(t time.Time) string {
	return t.Format(time.RFC3339)
}

// goodAllowedTimeDuration uses time.Duration which is allowlisted - no warning.
func goodAllowedTimeDuration(d time.Duration) int64 {
	return d.Nanoseconds()
}

// goodAllowedContext uses context.Context which is allowlisted - no warning.
func goodAllowedContext(ctx context.Context) error {
	return ctx.Err()
}

// goodNoParams is a function with no parameters - no warning.
func goodNoParams() string {
	return "no params"
}

// goodBuiltinTypes uses builtin types only - no warning.
func goodBuiltinTypes(s string, i int, b bool) string {
	// Vérification de la condition
	if b {
		return s
	}
	return string(rune(i))
}
