package test001

import "testing"

func TestBadPackageName(t *testing.T) {
	// Ce test utilise le mauvais nom de package
	t.Log("Devrait utiliser test001_test")
}
