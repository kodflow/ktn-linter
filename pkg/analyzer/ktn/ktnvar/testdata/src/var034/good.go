// Package var034 provides test cases for KTN-VAR-034.
package var034

import "sync"

// goodConcurrent utilise wg.Go() pour lancer des goroutines (Go 1.25+).
//
// Params:
//   - items: liste d'items a traiter
func goodConcurrent(items []int) {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	// Iteration sur les items avec wg.Go
	for _, item := range items {
		// Capture de la variable de boucle
		it := item
		wg.Go(func() {
			goodProcess(it)
		})
	}
	// Attente de la fin
	wg.Wait()
}

// goodSimple utilise wg.Go() pour une goroutine simple.
//
// Returns: nothing
func goodSimple() {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	wg.Go(func() {
		goodDoWork()
	})
	// Attente de la fin
	wg.Wait()
}

// goodMultiple utilise wg.Go() plusieurs fois.
//
// Returns: nothing
func goodMultiple() {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	wg.Go(func() {
		goodTask1()
	})
	wg.Go(func() {
		goodTask2()
	})
	// Attente de la fin
	wg.Wait()
}

// goodSeparateStatements a wg.Add(1) separe du go stmt.
//
// Returns: nothing
func goodSeparateStatements() {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)
	// Code entre Add et go
	x := 42
	_ = x
	go func() {
		defer wg.Done()
		goodDoWork()
	}()
	// Attente de la fin
	wg.Wait()
}

// goodAddNot1 utilise wg.Add avec une valeur differente de 1.
//
// Returns: nothing
func goodAddNot1() {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	// Add avec valeur > 1
	wg.Add(5)
	// Lancement de 5 goroutines
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			goodDoWork()
		}()
	}
	// Attente de la fin
	wg.Wait()
}

// goodNoDeferDone n'a pas defer wg.Done() au debut.
//
// Returns: nothing
func goodNoDeferDone() {
	// Declaration du WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// Pas de defer Done au debut
		goodDoWork()
		wg.Done()
	}()
	// Attente de la fin
	wg.Wait()
}

// goodDifferentWaitGroup utilise des WaitGroups differents.
//
// Returns: nothing
func goodDifferentWaitGroup() {
	// Declaration des WaitGroups
	var wg1, wg2 sync.WaitGroup
	wg1.Add(1)
	go func() {
		// Done sur wg2, pas wg1
		defer wg2.Done()
		goodDoWork()
	}()
	// Attente de la fin
	wg1.Wait()
	wg2.Wait()
}

// goodProcess traite un item.
//
// Params:
//   - it: item a traiter
func goodProcess(it int) {
	// Implementation
	_ = it
}

// goodDoWork effectue un travail.
func goodDoWork() {
	// Implementation
}

// goodTask1 effectue la tache 1.
func goodTask1() {
	// Implementation
}

// goodTask2 effectue la tache 2.
func goodTask2() {
	// Implementation
}
