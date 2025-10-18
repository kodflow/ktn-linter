package rules_func

// bad_function_name viole plusieurs règles simultanément.
func bad_function_name(a, b, c, d, e, f int) int {
	result := 0
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			if i%3 == 0 {
				if i%5 == 0 {
					if i%7 == 0 {
						result += a + b + c + d + e + f
					}
				}
			}
		}
	}
	// Early return from function.
	return result
}
