package var018

// DataHolder contient des données en bytes.
type DataHolder struct {
	content []byte
}

// GetBytes retourne le contenu en bytes.
//
// Returns:
//   - []byte: contenu
func (d *DataHolder) GetBytes() []byte {
	// Retour du contenu
	return d.content
}

// badRepeatedConversionInLoop convertit plusieurs fois dans une boucle.
//
// Params:
//   - data: données à parcourir
//   - target: valeur cible recherchée
//
// Returns:
//   - int: nombre d'occurrences
func badRepeatedConversionInLoop(data [][]byte, target string) int {
	count := 0
	// Parcours des données
	for _, item := range data {
		// Vérification de correspondance - CONVERSION RÉPÉTÉE
		if string(item) == target {
			count++
		}
	}
	// Retour du compteur
	return count
}

// badMultipleConversionsInFunction convertit plusieurs fois la même variable.
//
// Params:
//   - data: données à analyser
func badMultipleConversionsInFunction(data []byte) {
	// Vérification hello - CONVERSION #1
	if string(data) == "hello" {
		println("found hello")
	}
	// Vérification world - CONVERSION #2
	if string(data) == "world" {
		println("found world")
	}
	// Affichage - CONVERSION #3
	println(string(data))
}

// badConversionInForLoop convertit dans un for classique.
//
// Params:
//   - items: éléments à parcourir
func badConversionInForLoop(items [][]byte) {
	// Parcours des éléments
	for i := 0; i < len(items); i++ {
		// Vérification - CONVERSION RÉPÉTÉE
		if string(items[i]) == "test" {
			println("found")
		}
	}
}

// badNestedLoopConversion convertit dans une boucle imbriquée.
//
// Params:
//   - row: ligne de données
func badNestedLoopConversion(row [][]byte) {
	// Parcours des cellules
	for _, cell := range row {
		// Vérification - CONVERSION RÉPÉTÉE
		if string(cell) == "x" {
			println("found x")
		}
	}
}

// badMapKeyConversion utilise string() répété pour les clés de map.
//
// Params:
//   - cache: map de cache
//   - keys: clés à chercher
func badMapKeyConversion(cache map[string]int, keys [][]byte) {
	// Parcours des clés
	for _, key := range keys {
		// Accès à la map - CONVERSION RÉPÉTÉE
		_ = cache[string(key)]
	}
}

// badNestedIndexConversion teste la conversion avec index imbriqué.
//
// Params:
//   - matrix: matrice de données
func badNestedIndexConversion(matrix [][][]byte) {
	// Parcours des lignes
	for _, row := range matrix {
		// Parcours des cellules
		for _, cell := range row {
			// Vérification - CONVERSION RÉPÉTÉE
			if string(cell) == "test" {
				println("found")
			}
		}
	}
}

// badStructFieldConversion teste la conversion d'un champ de struct.
//
// Params:
//   - holders: liste de conteneurs
func badStructFieldConversion(holders []DataHolder) {
	// Parcours des conteneurs
	for _, h := range holders {
		// Vérification - CONVERSION RÉPÉTÉE
		if string(h.content) == "test" {
			println("found")
		}
	}
}

// badMethodResultConversion teste la conversion du résultat d'une méthode.
//
// Params:
//   - holders: liste de conteneurs
func badMethodResultConversion(holders []*DataHolder) {
	// Parcours des conteneurs
	for _, h := range holders {
		// Vérification avec appel de méthode - CONVERSION RÉPÉTÉE
		if string(h.GetBytes()) == "test" {
			println("found")
		}
	}
}

// badActualNestedIndex teste l'index imbriqué dans la conversion elle-même.
//
// Params:
//   - matrix: matrice de données
func badActualNestedIndex(matrix [][][]byte) {
	// Boucle sur la matrice
	for i := 0; i < len(matrix); i++ {
		// Boucle sur les lignes
		for j := 0; j < len(matrix[i]); j++ {
			// Conversion avec double index - CONVERSION RÉPÉTÉE
			if string(matrix[i][j]) == "test" {
				println("found")
			}
		}
	}
}

// badSelectorConversion teste la conversion avec sélecteur de struct.
//
// Params:
//   - items: liste d'éléments
func badSelectorConversion(items []DataHolder) {
	// Parcours des éléments
	for _, item := range items {
		// Vérification a - CONVERSION #1
		if string(item.content) == "a" {
			println("a")
		}
		// Vérification b - CONVERSION #2
		if string(item.content) == "b" {
			println("b")
		}
		// Vérification c - CONVERSION #3
		if string(item.content) == "c" {
			println("c")
		}
	}
}

// helperFunc retourne des données de test.
//
// Returns:
//   - []byte: données de test
func helperFunc() []byte {
	// Retour de données
	return []byte("test")
}

// badCallExprConversion teste avec un appel de fonction.
func badCallExprConversion() {
	// Vérification a - CONVERSION #1
	if string(helperFunc()) == "a" {
		println("a")
	}
	// Vérification b - CONVERSION #2
	if string(helperFunc()) == "b" {
		println("b")
	}
	// Vérification c - CONVERSION #3
	if string(helperFunc()) == "c" {
		println("c")
	}
}
