// Package var012 contains test cases for KTN rules.
package var012

// dataHolder contient des données en bytes.
// Cette structure est utilisée pour tester les conversions répétées.
type dataHolder struct {
	content []byte
}

// newDataHolder crée un nouveau dataHolder.
//
// Params:
//   - data: données initiales
//
// Returns:
//   - *dataHolder: nouveau dataHolder
func newDataHolder(data []byte) *dataHolder {
	// Retour d'une nouvelle instance
	return &dataHolder{content: data}
}

// bytes retourne le contenu en bytes.
//
// Returns:
//   - []byte: contenu
func (d *dataHolder) bytes() []byte {
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
	for range items {
		// Vérification - CONVERSION RÉPÉTÉE
		if string(items[0]) == "test" {
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
func badStructFieldConversion(holders []dataHolder) {
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
func badMethodResultConversion(holders []*dataHolder) {
	// Parcours des conteneurs
	for _, h := range holders {
		// Vérification avec appel de méthode - CONVERSION RÉPÉTÉE
		if string(h.bytes()) == "test" {
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
	for range matrix {
		// Boucle sur les lignes
		for range matrix[0] {
			// Conversion avec double index - CONVERSION RÉPÉTÉE
			if string(matrix[0][0]) == "test" {
				println("found")
			}
		}
	}
}

// badSelectorConversion teste la conversion avec sélecteur de struct.
//
// Params:
//   - items: liste d'éléments
func badSelectorConversion(items []dataHolder) {
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

// init utilise les fonctions privées
func init() {
	// Appel de newDataHolder
	_ = newDataHolder(nil)
	// Appel de badRepeatedConversionInLoop
	_ = badRepeatedConversionInLoop(nil, "")
	// Appel de badMultipleConversionsInFunction
	badMultipleConversionsInFunction(nil)
	// Appel de badConversionInForLoop
	badConversionInForLoop(nil)
	// Appel de badNestedLoopConversion
	badNestedLoopConversion(nil)
	// Appel de badMapKeyConversion
	badMapKeyConversion(nil, nil)
	// Appel de badNestedIndexConversion
	badNestedIndexConversion(nil)
	// Appel de badStructFieldConversion
	badStructFieldConversion(nil)
	// Appel de badMethodResultConversion
	badMethodResultConversion(nil)
	// Appel de badActualNestedIndex
	badActualNestedIndex(nil)
	// Appel de badSelectorConversion
	badSelectorConversion(nil)
	// Appel de helperFunc
	helperFunc()
	// Appel de badCallExprConversion
	badCallExprConversion()
}
