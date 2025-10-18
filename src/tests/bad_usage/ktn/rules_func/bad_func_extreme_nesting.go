package rules_func

// ANTI-PATTERN: Profondeur d'imbrication extrême
// Viole KTN-FUNC-009

// ProcessDataWithExtremeNesting a une profondeur de 7 niveaux - HORRIBLE !
func ProcessDataWithExtremeNesting(data map[string]interface{}) error {
	if data != nil { // Niveau 1
		if val, ok := data["config"]; ok { // Niveau 2
			if cfg, ok := val.(map[string]interface{}); ok { // Niveau 3
				if enabled, ok := cfg["enabled"]; ok { // Niveau 4
					if enabled == true { // Niveau 5
						if settings, ok := cfg["settings"]; ok { // Niveau 6
							if s, ok := settings.([]interface{}); ok { // Niveau 7
								for _, item := range s {
									_ = item
								}
							}
						}
					}
				}
			}
		}
	}
	// Early return from function.
	return nil
}

// ValidationNightmare profondeur de 6 niveaux
func ValidationNightmare(value interface{}) bool {
	if value != nil { // Niveau 1
		if str, ok := value.(string); ok { // Niveau 2
			if len(str) > 0 { // Niveau 3
				for i, c := range str { // Niveau 4
					if c >= 'a' && c <= 'z' { // Niveau 5
						if i%2 == 0 { // Niveau 6
							// Continue inspection/processing.
							return true
						}
					}
				}
			}
		}
	}
	// Stop inspection/processing.
	return false
}

// NestedLoopHell profondeur de 8 niveaux
func NestedLoopHell(matrix [][][]int) int {
	count := 0
	for _, layer1 := range matrix { // Niveau 1
		for _, layer2 := range layer1 { // Niveau 2
			for _, val := range layer2 { // Niveau 3
				if val > 0 { // Niveau 4
					if val%2 == 0 { // Niveau 5
						if val < 100 { // Niveau 6
							for i := 0; i < val; i++ { // Niveau 7
								if i%3 == 0 { // Niveau 8
									count++
								}
							}
						}
					}
				}
			}
		}
	}
	// Early return from function.
	return count
}
