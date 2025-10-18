package func005

// goodOneParam a un seul paramètre.
func goodOneParam(x int) {}

// goodThreeParams a trois paramètres.
func goodThreeParams(x, y, z int) {}

// goodFiveParams a exactement 5 paramètres (limite).
func goodFiveParams(a, b, c, d, e int) {}

// badSixParams a trop de paramètres.
func badSixParams(a, b, c, d, e, f int) {} // want "KTN-FUNC-005.*trop de paramètres.*6 > 5"

// badSevenParams a bien trop de paramètres.
func badSevenParams(a, b, c, d, e, f, g int) {} // want "KTN-FUNC-005.*trop de paramètres.*7 > 5"

// badManyParams a beaucoup trop de paramètres.
func badManyParams(a, b, c, d, e, f, g, h, i, j int) {} // want "KTN-FUNC-005.*trop de paramètres.*10 > 5"
