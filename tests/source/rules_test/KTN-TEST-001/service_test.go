// ❌ CAS INCORRECT: Package name sans suffixe _test (devrait être KTN_TEST_001_test)
package KTN_TEST_001

import "testing"

// TestAddUser teste l'ajout d'utilisateur.
func TestAddUser(t *testing.T) {
	service := NewUserService()
	service.AddUser("1", "Alice")

	name, exists := service.GetUser("1")
	if !exists {
		t.Fatal("user not found")
	}
	if name != "Alice" {
		t.Errorf("expected Alice, got %s", name)
	}
}
