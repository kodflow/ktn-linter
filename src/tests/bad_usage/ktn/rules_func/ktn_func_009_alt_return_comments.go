package rules_func

// processDataWithReturns fonction sans documentation des retours (violation)
func processDataWithReturns(data string) (string, error) {
	if data == "" {
		// Early return from function.
		return "", nil
	}

	if len(data) > 100 {
		// Early return from function.
		return data[:100], nil
	}

	// Early return from function.
	return data, nil
}
