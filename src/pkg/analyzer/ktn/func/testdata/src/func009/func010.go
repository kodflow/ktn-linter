package func010

// goodShallow a peu d'imbrication.
func goodShallow(a, b bool) int {
	if a {
		if b {
			return 1
		}
	}
	return 0
}

// goodThreeLevels a 3 niveaux (limite acceptable).
func goodThreeLevels(a, b, c bool) int {
	if a {
		if b {
			if c {
				return 1
			}
		}
	}
	return 0
}

// badFourLevels a 4 niveaux (trop profond).
func badFourLevels(a, b, c, d bool) int { // want "KTN-FUNC-010.*profondeur d'imbrication trop élevée"
	if a {
		if b {
			if c {
				if d {
					return 1
				}
			}
		}
	}
	return 0
}

// badWithLoops a trop d'imbrication avec boucles.
func badWithLoops(items [][]int) int { // want "KTN-FUNC-010.*profondeur d'imbrication trop élevée"
	for _, row := range items {
		for _, val := range row {
			if val > 0 {
				for i := 0; i < val; i++ {
					return i
				}
			}
		}
	}
	return 0
}

// badWithSwitch a trop d'imbrication avec switch.
func badWithSwitch(x, y int) int { // want "KTN-FUNC-010.*profondeur d'imbrication trop élevée"
	if x > 0 {
		switch y {
		case 1:
			if x == y {
				for i := 0; i < x; i++ {
					return i
				}
			}
		}
	}
	return 0
}
