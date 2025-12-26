// Mock API for testing.
package mock_only

// MockAPI is a mock implementation.
type MockAPI struct{}

// Execute runs the mock API.
//
// Params:
//   - request: API request
//
// Returns:
//   - string: API response
func (m *MockAPI) Execute(request string) string {
	// Return mock response
	return "response:" + request
}
