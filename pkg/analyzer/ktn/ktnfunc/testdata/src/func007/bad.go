// Package func007 contains test cases for KTN rules.
package func007

const (
	// CACHE_MULTIPLIER est le facteur de multiplication pour le cache
	CACHE_MULTIPLIER int = 2
	// DEFAULT_CACHE_VALUE est la valeur par défaut du cache
	DEFAULT_CACHE_VALUE int = 100
	// FIRST_ELEMENT_VALUE est la valeur du premier élément
	FIRST_ELEMENT_VALUE int = 100
	// ITEM_VALUE est la valeur de l'item
	ITEM_VALUE int = 200
)

// CounterData contient un compteur et des données.
// Utilisé pour tester les violations de getters avec side effects.
type CounterData struct {
	count int
	cache int
	data  []int
	items map[string]int
}

// CounterDataInterface définit les méthodes publiques de CounterData.
type CounterDataInterface interface {
	IsCountPositive() bool
	IsCacheReady() bool
	IsReady() bool
	HasData() bool
	IsFirstElementSet() bool
	HasItem() bool
}

// NewCounterData crée une nouvelle instance de CounterData.
//
// Returns:
//   - *CounterData: nouvelle instance
func NewCounterData() *CounterData {
	// Retour de la nouvelle instance
	return &CounterData{}
}

// IsCountPositive vérifie si le compteur est positif (side effect: modification de count)
//
// Returns:
//   - bool: true si le compteur est positif après incrémentation
func (c *CounterData) IsCountPositive() bool {
	// Incrémentation du compteur (side effect volontaire pour bad.go)
	c.count++
	// Retourne vrai si le compteur est positif
	return c.count > 0
}

// IsCacheReady vérifie si le cache est prêt (side effect: modification de cache)
//
// Returns:
//   - bool: true si le cache est positif
func (c *CounterData) IsCacheReady() bool {
	// Calcul et assignation du cache (side effect volontaire pour bad.go)
	c.cache = c.count * CACHE_MULTIPLIER
	// Retourne vrai si le cache est positif
	return c.cache > 0
}

// IsReady vérifie si le compteur est prêt (side effect: modification de cache)
//
// Returns:
//   - bool: true si le compteur est positif
func (c *CounterData) IsReady() bool {
	// Assignation du cache (side effect volontaire pour bad.go)
	c.cache = DEFAULT_CACHE_VALUE
	// Retourne vrai si le compteur est positif
	return c.count > 0
}

// HasData vérifie si le compteur contient des données (side effect: modification de count)
//
// Returns:
//   - bool: true si le compteur est positif après incrémentation
func (c *CounterData) HasData() bool {
	// Incrémentation du compteur (side effect volontaire pour bad.go)
	c.count++
	// Retourne vrai si le compteur est positif
	return c.count > 0
}

// IsFirstElementSet vérifie si le premier élément est défini (side effect: modification du slice)
//
// Returns:
//   - bool: true si le premier élément est défini
func (c *CounterData) IsFirstElementSet() bool {
	// Modification du premier élément (side effect volontaire pour bad.go)
	c.data[0] = FIRST_ELEMENT_VALUE
	// Retourne vrai si le premier élément est défini
	return c.data[0] > 0
}

// HasItem vérifie si un item existe (side effect: modification de la map)
//
// Returns:
//   - bool: true si l'item existe
func (c *CounterData) HasItem() bool {
	// Modification de l'item (side effect volontaire pour bad.go)
	c.items["key"] = ITEM_VALUE
	// Retourne vrai si l'item existe
	return c.items["key"] > 0
}
