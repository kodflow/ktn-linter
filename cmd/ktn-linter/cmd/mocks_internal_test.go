// Internal test mocks for cmd package.
package cmd

import "github.com/kodflow/ktn-linter/pkg/updater"

// mockFlagGetter is a mock implementation of flagGetter for testing.
type mockFlagGetter struct {
	boolValues   map[string]bool
	stringValues map[string]string
	boolErr      error
	stringErr    error
}

// newMockFlagGetter creates a new mock flag getter.
func newMockFlagGetter() *mockFlagGetter {
	return &mockFlagGetter{
		boolValues:   make(map[string]bool),
		stringValues: make(map[string]string),
	}
}

// GetBool returns a mock bool value.
func (m *mockFlagGetter) GetBool(name string) (bool, error) {
	// Check for error
	if m.boolErr != nil {
		return false, m.boolErr
	}
	return m.boolValues[name], nil
}

// GetString returns a mock string value.
func (m *mockFlagGetter) GetString(name string) (string, error) {
	// Check for error
	if m.stringErr != nil {
		return "", m.stringErr
	}
	return m.stringValues[name], nil
}

// mockUpdaterService is a mock implementation of updaterService for testing.
type mockUpdaterService struct {
	checkInfo   updater.UpdateInfo
	checkErr    error
	upgradeInfo updater.UpdateInfo
	upgradeErr  error
}

// CheckForUpdate returns mock update info.
func (m *mockUpdaterService) CheckForUpdate() (updater.UpdateInfo, error) {
	return m.checkInfo, m.checkErr
}

// Upgrade returns mock upgrade info.
func (m *mockUpdaterService) Upgrade() (updater.UpdateInfo, error) {
	return m.upgradeInfo, m.upgradeErr
}
