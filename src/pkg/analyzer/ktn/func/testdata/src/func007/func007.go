package func007

// goodSimple a une complexité acceptable.
func goodSimple(x int) int {
	if x > 0 {
		return x * 2
	}
	if x < 0 {
		return -x
	}
	return 0
}

// goodLimit est à la limite de complexité 10.
func goodLimit(x int) int {
	if x == 1 {
		return 1
	}
	if x == 2 {
		return 2
	}
	if x == 3 {
		return 3
	}
	if x == 4 {
		return 4
	}
	if x == 5 {
		return 5
	}
	if x == 6 {
		return 6
	}
	if x == 7 {
		return 7
	}
	if x == 8 {
		return 8
	}
	if x == 9 {
		return 9
	}
	return 0
}

// badTooComplex a une complexité trop élevée.
func badTooComplex(x int) int { // want "KTN-FUNC-007.*complexité cyclomatique trop élevée"
	if x > 0 {
		if x > 1 {
			if x > 2 {
				if x > 3 {
					if x > 4 {
						if x > 5 {
							if x > 6 {
								if x > 7 {
									if x > 8 {
										if x > 9 {
											if x > 10 {
												return x
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return 0
}

// badComplexLogic a trop d'opérateurs logiques.
func badComplexLogic(a, b, c, d, e, f bool) bool { // want "KTN-FUNC-007.*complexité cyclomatique trop élevée"
	if a && b || c && d {
		if e && f {
			if a || b {
				if c || d {
					if e || f {
						if a && c {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
