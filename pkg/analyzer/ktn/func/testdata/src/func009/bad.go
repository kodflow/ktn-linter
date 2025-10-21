package func009

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

// Counter est un compteur avec cache
type Counter struct {
	count int
	cache int
}

// GetCount retourne le compteur incrémenté (side effect: modification de count)
//
// Returns:
//   - int: la valeur du compteur après incrémentation
func (c *Counter) GetCount() int {
	// Incrémentation du compteur (side effect volontaire pour bad.go)
	c.count++
	// Retourne le compteur incrémenté
	return c.count
}

// GetCachedValue retourne la valeur en cache calculée (side effect: modification de cache)
//
// Returns:
//   - int: la valeur en cache calculée
func (c *Counter) GetCachedValue() int {
	// Calcul et assignation du cache (side effect volontaire pour bad.go)
	c.cache = c.count * CACHE_MULTIPLIER
	// Retourne la valeur en cache
	return c.cache
}

// IsReady vérifie si le compteur est prêt (side effect: modification de cache)
//
// Returns:
//   - bool: true si le compteur est positif
func (c *Counter) IsReady() bool {
	// Assignation du cache (side effect volontaire pour bad.go)
	c.cache = DEFAULT_CACHE_VALUE
	// Retourne vrai si le compteur est positif
	return c.count > 0
}

// HasData vérifie si le compteur contient des données (side effect: modification de count)
//
// Returns:
//   - bool: true si le compteur est positif après incrémentation
func (c *Counter) HasData() bool {
	// Incrémentation du compteur (side effect volontaire pour bad.go)
	c.count++
	// Retourne vrai si le compteur est positif
	return c.count > 0
}

// DataStore stocke des données dans un slice et une map
type DataStore struct {
	data  []int
	items map[string]int
}

// GetFirstElement retourne le premier élément (side effect: modification du slice)
//
// Returns:
//   - int: la valeur du premier élément
func (d *DataStore) GetFirstElement() int {
	// Modification du premier élément (side effect volontaire pour bad.go)
	d.data[0] = FIRST_ELEMENT_VALUE
	// Retourne le premier élément
	return d.data[0]
}

// GetItem retourne un item de la map (side effect: modification de la map)
//
// Returns:
//   - int: la valeur de l'item
func (d *DataStore) GetItem() int {
	// Modification de l'item (side effect volontaire pour bad.go)
	d.items["key"] = ITEM_VALUE
	// Retourne l'item
	return d.items["key"]
}
