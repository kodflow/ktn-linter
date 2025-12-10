package good

// Private helper function that uses underscore test naming convention
func populateData(data string) string {
	return "populated: " + data
}

// Another private helper
func validateInput(input string) bool {
	return len(input) > 0
}
