// ✅ CORRIGÉ: helper.go existe (pas de fichier _test.go orphelin)
package KTN_TEST_003_GOOD_test

import "testing"

// TestFormat teste le formatage.
//
// Params:
//   - t: instance de test
func TestFormat(t *testing.T) {
	h := &Helper{}
	result := h.Format("test")
	if result != "[test]" {
		t.Errorf("expected [test], got %s", result)
	}
}
