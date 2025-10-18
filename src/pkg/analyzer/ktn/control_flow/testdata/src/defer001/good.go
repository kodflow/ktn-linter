package defer001

import "os"

// Cas corrects - defer utilisé correctement

// GoodDeferInSeparateFunc - pattern recommandé avec fonction séparée
func GoodDeferInSeparateFunc(files []string) {
	for _, file := range files {
		processFile(file)
	}
}

func processFile(name string) {
	f, _ := os.Open(name)
	defer f.Close() // OK - dans fonction séparée
	// traitement...
}

// GoodMultipleDefer - plusieurs defer hors boucle
func GoodMultipleDefer() {
	f1, _ := os.Open("file1.txt")
	defer f1.Close()

	f2, _ := os.Open("file2.txt")
	defer f2.Close()
}

// GoodDeferAfterLoop - defer après la boucle
func GoodDeferAfterLoop() {
	items := []int{1, 2, 3}
	for range items {
		// traitement
	}

	f, _ := os.Open("result.txt")
	defer f.Close()
}

// GoodNoDefer - pas de defer du tout
func GoodNoDefer(x int) int {
	return x * 2
}
