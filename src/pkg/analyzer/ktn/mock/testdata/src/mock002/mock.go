package mock002

// Mock incomplet - manque MockDataStore
// MockUserService represents the struct.
type MockUserService struct {
	GetUserFunc func(id int) (string, error)
}

func (m *MockUserService) GetUser(id int) (string, error) {
	return m.GetUserFunc(id)
}
