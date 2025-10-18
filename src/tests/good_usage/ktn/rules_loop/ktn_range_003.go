package rules_loop

import (
	"fmt"
	"sync"
)

// ✅ GOOD: copie locale avant goroutine AVEC synchronisation
func processAsyncGood(items []string) {
	var wg sync.WaitGroup
	for _, item := range items {
		item := item // ✅ copie locale (shadow)
		wg.Add(1)
		go asyncWorker(item, &wg)
	}
	wg.Wait() // ✅ synchronisation
}

// ✅ GOOD: passage par paramètre AVEC synchronisation
func processAsyncParam(items []string) {
	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go asyncWorker(item, &wg)
	}
	wg.Wait() // ✅ synchronisation
}

// ✅ GOOD: defer avec copie locale dans closure pour éviter accumulation
func deferInLoopGood(files []string) {
	for _, filename := range files {
		processDeferFile(filename)
	}
}

// processDeferFile traite un fichier avec defer dans closure.
//
// Params:
//   - fn: nom du fichier
func processDeferFile(fn string) {
	defer func() {
		fmt.Println("Processed:", fn)
	}()
	// Traitement du fichier ici (defer s'exécute à la fin de cette fonction)
}

// ✅ GOOD: goroutine avec channel et copie AVEC synchronisation
func fanOutGood(items []int, results chan int) {
	var wg sync.WaitGroup
	for _, v := range items {
		v := v // ✅ copie locale
		wg.Add(1)
		go fanOutWorker(v, results, &wg)
	}
	wg.Wait() // ✅ synchronisation
}

// ✅ GOOD: les deux variables copiées AVEC synchronisation
func capturesBothGood(items []string) {
	var wg sync.WaitGroup
	for i, item := range items {
		i := i       // ✅ copie index
		item := item // ✅ copie item
		wg.Add(1)
		go printWorker(i, item, &wg)
	}
	wg.Wait() // ✅ synchronisation
}

// ✅ GOOD: nested avec copies AVEC synchronisation
func nestedClosuresGood(matrix [][]int) {
	var outerWg sync.WaitGroup
	for _, row := range matrix {
		row := row // ✅ copie outer
		outerWg.Add(1)
		go nestedOuterWorker(row, &outerWg)
	}
	outerWg.Wait() // ✅ synchronisation outer
}

// ✅ GOOD: function literals avec copie
func assignFunctionsGood(items []string) []func() {
	var funcs []func()
	for _, item := range items {
		item := item // ✅ copie pour closure
		funcs = append(funcs, func() {
			processFuncGood(item) // ✅ chaque fonction a sa valeur
		})
	}
	// Early return from function.
	return funcs
}

// ✅ GOOD: map avec copies AVEC synchronisation
func processMapGoodAsync(m map[string]int) {
	var wg sync.WaitGroup
	for key, value := range m {
		key := key     // ✅ copie key
		value := value // ✅ copie value
		wg.Add(1)
		go mapWorker(key, value, &wg)
	}
	wg.Wait() // ✅ synchronisation
}

// ✅ GOOD: channel avec paramètre
func sendToChannelGood(items []string, ch chan func()) {
	for _, item := range items {
		item := item
		ch <- func() {
			processChannelGood(item) // ✅ copie
		}
	}
}

// ✅ GOOD: alternative - collect puis process
func processInTwoPhases(items []string) {
	// Phase 1: collect work
	type work struct{ item string }
	works := make([]work, len(items))
	for i, item := range items {
		works[i] = work{item: item} // ✅ copié dans struct
	}

	// Phase 2: process async AVEC synchronisation
	var wg sync.WaitGroup
	for _, w := range works {
		w := w
		wg.Add(1)
		go workWorker(w, &wg)
	}
	wg.Wait() // ✅ synchronisation
}

// ✅ GOOD: utiliser index pour accès AVEC synchronisation
func processWithIndex(items []string) {
	var wg sync.WaitGroup
	for i := range items {
		i := i // ✅ copie index
		wg.Add(1)
		go indexWorker(i, items, &wg)
	}
	wg.Wait() // ✅ synchronisation
}

// Fonctions helper

// asyncWorker traite un élément de manière asynchrone.
//
// Params:
//   - item: élément à traiter
//   - wg: waitgroup
func asyncWorker(item string, wg *sync.WaitGroup) {
	defer wg.Done()
	processAsyncGood2(item)
}

// fanOutWorker traite une valeur et envoie le résultat au channel.
//
// Params:
//   - v: valeur à traiter
//   - results: channel de résultats
//   - wg: waitgroup
func fanOutWorker(v int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	results <- computeGood(v)
}

// printWorker affiche l'index et l'élément.
//
// Params:
//   - i: index
//   - item: élément
//   - wg: waitgroup
func printWorker(i int, item string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%d: %s\n", i, item)
}

// nestedOuterWorker traite une ligne de la matrice.
//
// Params:
//   - row: ligne de la matrice
//   - outerWg: waitgroup externe
func nestedOuterWorker(row []int, outerWg *sync.WaitGroup) {
	defer outerWg.Done()
	var innerWg sync.WaitGroup
	for _, val := range row {
		val := val
		innerWg.Add(1)
		go nestedInnerWorker(val, &innerWg)
	}
	innerWg.Wait()
}

// nestedInnerWorker traite une valeur de la matrice.
//
// Params:
//   - val: valeur à traiter
//   - innerWg: waitgroup interne
func nestedInnerWorker(val int, innerWg *sync.WaitGroup) {
	defer innerWg.Done()
	processNestedGood(val)
}

// mapWorker traite une paire clé-valeur du map.
//
// Params:
//   - key: clé
//   - value: valeur
//   - wg: waitgroup
func mapWorker(key string, value int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s: %d\n", key, value)
}

// workWorker traite un work item.
//
// Params:
//   - w: work item
//   - wg: waitgroup
func workWorker(w interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	processWork(w)
}

// indexWorker traite un élément par son index.
//
// Params:
//   - i: index
//   - items: slice d'éléments
//   - wg: waitgroup
func indexWorker(i int, items []string, wg *sync.WaitGroup) {
	defer wg.Done()
	processAsyncGood2(items[i])
}

func processAsyncGood2(s string)  {}
func computeGood(v int) int       { return v * 2 }
func processNestedGood(v int)     {}
func processFuncGood(s string)    {}
func processChannelGood(s string) {}
func processWork(w interface{})   {}
