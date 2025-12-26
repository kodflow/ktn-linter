package bad

// Private function without test - should trigger KTN-TEST-002
func formatMessage(msg string) string {
	return "[INFO] " + msg
}
