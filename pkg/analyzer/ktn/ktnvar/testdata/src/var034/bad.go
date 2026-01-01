// Package var034 provides test cases for KTN-VAR-034.
package var034

import "sync"

// badConcurrent utilise le pattern wg.Add(1) + go func() + defer wg.Done()
// au lieu de wg.Go() disponible depuis Go 1.25.
//
// Returns: nothing
func badConcurrent(items []int) {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	// Iteration sur les items
	for _, item := range items {
		wg.Add(1)
		go func(it int) { // want "KTN-VAR-034"
			defer wg.Done()
			process(it)
		}(item)
	}
	// Attente de la fin
	wg.Wait()
}

// badSimple utilise le pattern wg.Add(1) + go func() + defer wg.Done()
// pour une goroutine simple.
//
// Returns: nothing
func badSimple() {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { // want "KTN-VAR-034"
		defer wg.Done()
		doWork()
	}()
	// Attente de la fin
	wg.Wait()
}

// badMultiple utilise le pattern multiple fois.
//
// Returns: nothing
func badMultiple() {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { // want "KTN-VAR-034"
		defer wg.Done()
		task1()
	}()
	wg.Add(1)
	go func() { // want "KTN-VAR-034"
		defer wg.Done()
		task2()
	}()
	// Attente de la fin
	wg.Wait()
}

// process traite un item.
//
// Params:
//   - it: item a traiter
func process(it int) {
	// Implementation
	_ = it
}

// doWork effectue un travail.
func doWork() {
	// Implementation
}

// task1 effectue la tache 1.
func task1() {
	// Implementation
}

// task2 effectue la tache 2.
func task2() {
	// Implementation
}
