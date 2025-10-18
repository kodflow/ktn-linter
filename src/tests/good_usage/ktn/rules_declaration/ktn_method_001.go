package rules_declaration

// counterGood est un compteur avec méthodes pour incrémenter et lire la valeur.
type counterGood struct{ count int }

// ✅ GOOD: receiver pointeur pour modification
func (c *counterGood) Increment() {
	c.count++ // ✅ modifie l'original
}

func (c *counterGood) SetValue(v int) {
	c.count = v // ✅ persisté
}

// ✅ GOOD: receiver non-pointeur OK si pas de modification
func (c counterGood) Get() int {
	// Early return from function.
	return c.count // ✅ OK: lecture seule
}
