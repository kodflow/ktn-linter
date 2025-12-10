package bad

// Private function without test - should trigger KTN-TEST-003
func formatMessage(msg string) string {
	return "[INFO] " + msg
}
