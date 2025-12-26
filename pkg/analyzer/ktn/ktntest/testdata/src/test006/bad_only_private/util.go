package test008 // want "KTN-TEST-006: le fichier 'util.go' contient des fonctions privées. Il doit avoir un fichier 'util_internal_test.go' \\(white-box\\)"

// parseValue parse une valeur (fonction privée)
func parseValue(input string) int {
	return len(input)
}

// formatData formate les données (fonction privée)
func formatData(data string) string {
	return "[" + data + "]"
}
