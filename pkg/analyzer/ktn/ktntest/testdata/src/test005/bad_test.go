package test007

import "testing"

// TestBadWithSkip utilise t.Skip()
func TestBadWithSkip(t *testing.T) {
	t.Skip("Test skipped for some reason") // want "KTN-TEST-005: t.Skip\\(\\) interdit dans 't.Skip\\(\\)'. Les tests doivent passer"
	// Rest of test never runs
}

// TestBadWithSkipf utilise t.Skipf()
func TestBadWithSkipf(t *testing.T) {
	t.Skipf("Test skipped: %s", "some reason") // want "KTN-TEST-005: t.Skip\\(\\) interdit dans 't.Skipf\\(\\)'. Les tests doivent passer"
}

// TestBadWithSkipNow utilise t.SkipNow()
func TestBadWithSkipNow(t *testing.T) {
	t.SkipNow() // want "KTN-TEST-005: t.Skip\\(\\) interdit dans 't.SkipNow\\(\\)'. Les tests doivent passer"
}
