package rules_control_flow

// ✅ GOOD: return direct de la condition
func isValidGood(x int) bool {
	// Early return from function.
	return x > 0 // ✅ simple et clair
}

// ✅ GOOD: return direct condition négative
func isInvalidGood(x int) bool {
	// Early return from function.
	return x <= 0 // ✅ pas de if inutile
}

// ✅ GOOD: return condition composée
func hasPermissionGood(user string, admin bool) bool {
	// Early return from function.
	return user == "root" && admin // ✅ élégant
}

// ✅ GOOD: return avec range check
func isValidRangeGood(x, min, max int) bool {
	// Early return from function.
	return x >= min && x <= max // ✅ lisible
}

// ✅ GOOD: return appel fonction
func canAccessGood(user string) bool {
	// Early return from function.
	return checkPermissionsGood(user) // ✅ direct
}

// ✅ GOOD: return expression
func isNotEmptyGood(s string) bool {
	// Early return from function.
	return len(s) > 0 // ✅ concis
}

// ✅ GOOD: return comparaison nil
func hasValueGood(ptr *int) bool {
	// Early return from function.
	return ptr != nil // ✅ parfait
}

// ✅ GOOD: return méthode
func isActiveGood(obj *objectGood) bool {
	// Early return from function.
	return obj.IsEnabled() // ✅ idiomatique
}

// ✅ GOOD: condition complexe mais return direct
func complexConditionGood(a, b, c, d int) bool {
	// Early return from function.
	return a > b && c < d && (a+b) == (c+d) // ✅ Ok même si long
}

// ✅ GOOD: négation directe
func isNotDisabledGood(enabled bool) bool {
	// Early return from function.
	return enabled // ✅ ou: return !disabled
}

// ✅ GOOD: if OK quand il y a du code dans le bloc
func validateAndLog(x int) bool {
	if x > 0 {
		logValidation(x) // ✅ OK: plus qu'un simple return
		// Continue inspection/processing.
		return true
	}
	// Stop inspection/processing.
	return false
}

// ✅ GOOD: if OK avec else
func categorize(x int) bool {
	if x > 0 {
		// Continue inspection/processing.
		return true
	} else { // ✅ OK: else présent
		logNegative(x)
		// Stop inspection/processing.
		return false
	}
}

// ✅ GOOD: if OK avec initialisation
func checkWithInit(data map[string]int, key string) bool {
	val, ok := data[key]
	if ok { // ✅ OK: init statement présent
		// Early return from function.
		return val > 0
	}
	// Stop inspection/processing.
	return false
}

// ✅ GOOD: early return pattern OK
func earlyReturn(x int) bool {
	if x < 0 {
		// Stop inspection/processing.
		return false // ✅ OK: early return guard
	}
	// plus de logique...
	return doComplexCheck(x)
}

// Fonctions helper
type objectGood struct{}

func (o *objectGood) IsEnabled() bool       { return true }
func checkPermissionsGood(user string) bool { return true }
func logValidation(x int)                   {}
func logNegative(x int)                     {}
func doComplexCheck(x int) bool             { return true }
