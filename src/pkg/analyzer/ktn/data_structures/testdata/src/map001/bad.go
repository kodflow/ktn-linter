package map001

// TODO: L'analyseur MAP-001 a des limitations
// Il ne détecte pas les écritures sur des paramètres de fonction
// car il considère que l'appelant a initialisé la map

func GoodMapWithCheck(m map[string]int) {
	if m != nil {
		m["key"] = 42
	}
}

func GoodMapInitialized() {
	m := make(map[string]int)
	m["key"] = 42
}
