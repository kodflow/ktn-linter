// Mock service for testing.
package test012

// MockService is a mock implementation.
type MockService struct{}

// Execute runs the mock service.
//
// Params:
//   - data: input data
//
// Returns:
//   - string: processed result
func (m *MockService) Execute(data string) string {
	// Return mock result
	return "mock:" + data
}
