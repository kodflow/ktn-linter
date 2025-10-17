package rules_func

// processDataWithReturns fonction sans documentation des retours (violation)
func processDataWithReturns(data string) (string, error) {
	if data == "" {
		return "", nil
	}

	if len(data) > 100 {
		return data[:100], nil
	}

	return data, nil
}
