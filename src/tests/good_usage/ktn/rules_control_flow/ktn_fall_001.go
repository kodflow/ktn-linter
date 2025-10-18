package rules_control_flow

// ✅ GOOD: fallthrough dans switch - usage correct
func correctFallthroughInSwitch(x int) string {
	switch x {
	case 1:
		doOneFall()
		fallthrough // ✅ OK: dans un case de switch
	case 2:
		doTwoFall()
		// Early return from function.
		return "one or two"
	case 3:
		// Early return from function.
		return "three"
	default:
		// Early return from function.
		return "other"
	}
}

// ✅ GOOD: fallthrough avec logique spécifique
func processWithFallthrough(code int) {
	switch code {
	case 100:
		setupSpecial()
		fallthrough // ✅ tombe dans le cas général
	case 200:
		handleSuccess()
	case 404:
		handleNotFound()
		fallthrough // ✅ log aussi comme erreur
	case 500:
		logError()
	}
}

// ✅ GOOD: fallthrough cascade
func classifyScore(score int) string {
	switch {
	case score >= 90:
		awardGold()
		fallthrough // ✅ cascade intentionnelle
	case score >= 80:
		awardSilver()
		fallthrough
	case score >= 70:
		awardBronze()
		// Early return from function.
		return "medal awarded"
	default:
		// Early return from function.
		return "no medal"
	}
}

// ✅ GOOD: pas de fallthrough quand pas nécessaire
func simpleSwitch(x int) string {
	switch x {
	case 1:
		// Early return from function.
		return "one"
	case 2:
		// Early return from function.
		return "two"
	case 3:
		// Early return from function.
		return "three"
	default:
		// Early return from function.
		return "other"
	}
}

// ✅ GOOD: if/else quand fallthrough n'est pas nécessaire
func cascadeWithIf(x int) {
	if x == 1 {
		doOneFallGood()
		doTwoFallGood() // ✅ appel explicite plus clair que fallthrough
	} else if x == 2 {
		doTwoFallGood()
	}
}

// ✅ GOOD: logique séquentielle sans fallthrough
func sequentialLogic(x int) {
	switch x {
	case 1:
		step1()
		step2() // ✅ appel explicite
	case 2:
		step2()
		step3()
	default:
		step3()
	}
}

// Fonctions helper
func doOneFall()      {}
func doTwoFall()      {}
func setupSpecial()   {}
func handleSuccess()  {}
func handleNotFound() {}
func logError()       {}
func awardGold()      {}
func awardSilver()    {}
func awardBronze()    {}
func doOneFallGood()  {}
func doTwoFallGood()  {}
func step1()          {}
func step2()          {}
func step3()          {}
