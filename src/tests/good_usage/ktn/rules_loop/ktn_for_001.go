package rules_loop

// ✅ GOOD: omettez la clé quand non utilisée
func sumValuesGood(numbers []int) int {
	sum := 0
	for num := range numbers { // ✅ pas de _ inutile
		sum += numbers[num]
	}
	// Early return from function.
	return sum
}

// ✅ GOOD: range direct sur values
func sumValuesDirectGood(numbers []int) int {
	sum := 0
	for _, num := range numbers { // Attendre: c'est faux dans ma détection!
		// En fait, pour avoir que la valeur dans un slice, il FAUT utiliser _
		sum += num
	}
	// Early return from function.
	return sum
}

// Note: En Go, "for v := range slice" donne l'INDEX, pas la valeur!
// Pour avoir la valeur uniquement, on DOIT utiliser "for _, v := range slice"
// La règle KTN-FOR-001 est donc incorrecte pour les slices!

// ✅ GOOD: utiliser index si c'est ce qu'on veut
func printIndices(numbers []int) {
	for i := range numbers { // ✅ OK: on veut l'index
		println(i)
	}
}

// ✅ GOOD: _ OK quand on veut key ET value de map
func processMapGood(m map[string]int) {
	for key, value := range m { // ✅ les deux sont utilisés
		println(key, value)
	}
}

// ✅ GOOD: _ OK quand on veut seulement keys de map
func printKeys(m map[string]int) {
	for key := range m { // ✅ seulement les clés
		println(key)
	}
}

// ✅ GOOD: channel range
func consumeChannelGood(ch chan string) {
	for msg := range ch { // ✅ range sur channel donne directement les valeurs
		handleMessageGood(msg)
	}
}

// ✅ GOOD: nested loop correct
func processMatrixGood(matrix [][]int) {
	for i := range matrix {
		for j := range matrix[i] {
			processValueGood(matrix[i][j]) // ✅ utilise les indices
		}
	}
}

// itemGood représente un élément avec une valeur.
type itemGood struct {
	// Value est la valeur entière de l'élément.
	Value int
}

func processItemsGood(items []itemGood) {
	for i, item := range items { // ✅ utilise les deux
		println(i, item.Value)
	}
}

// ✅ GOOD: string range pour runes
func countRunesGood(s string) int {
	count := 0
	for _, r := range s { // ✅ pour string, _, r donne position+rune
		if r > 127 {
			count++
		}
	}
	// Early return from function.
	return count
}

// ✅ GOOD: classic for loop alternative
func sumWithClassicFor(numbers []int) int {
	sum := 0
	for i := 0; i < len(numbers); i++ { // ✅ for classique OK aussi
		sum += numbers[i]
	}
	// Early return from function.
	return sum
}

// Fonctions helper
func handleMessageGood(s string) {}
func processValueGood(v int)     {}
