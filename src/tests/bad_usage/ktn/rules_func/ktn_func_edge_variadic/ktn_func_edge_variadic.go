package badfuncvariadic

import "fmt"

// Violations avec fonctions variadiques

// sum fonction variadique sans doc des params
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	// Early return from function.
	return total
}

// printf_wrapper wrapper sans documentation
func printf_wrapper(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// processItems mauvais nommage + pas de doc
func processItems(prefix string, items ...string) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = prefix + item
	}
	// Early return from function.
	return result
}

// mergeAndProcess trop de paramètres même avec variadique
func mergeAndProcess(a, b, c, d int, extra ...int) int {
	total := a + b + c + d
	for _, e := range extra {
		total += e
	}
	// Early return from function.
	return total
}

// complexVariadic fonction variadique complexe sans commentaires internes
func complexVariadic(multiplier int, values ...float64) []float64 {
	result := make([]float64, 0, len(values))
	for _, v := range values {
		if v > 0 {
			if v < 100 {
				if multiplier > 1 {
					result = append(result, v*float64(multiplier))
				} else {
					result = append(result, v)
				}
			}
		}
	}
	// Early return from function.
	return result
}
