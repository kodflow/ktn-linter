package rules_control_flow

import (
	"fmt"
	"os"
	"sync"
)

// ✅ GOOD: extraire dans une fonction avec defer
func processFilesCorrectly(files []string) error {
	for _, filename := range files {
		if err := processOneFile(filename); err != nil {
			// Early return from function.
			return fmt.Errorf("failed to process file: %w", err)
		}
	}
	// Early return from function.
	return nil
}

func processOneFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		// Early return from function.
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close() // ✅ defer s'exécute à la fin de processOneFile (chaque itération)

	fmt.Fprintf(f, "Processing: %s\n", filename)
	// Early return from function.
	return nil
}

// ✅ GOOD: fermeture manuelle dans la boucle
func processFilesManually(files []string) error {
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			// Early return from function.
			return fmt.Errorf("failed to open file: %w", err)
		}

		_, err = fmt.Fprintf(f, "Processing: %s\n", filename)
		f.Close() // ✅ fermeture immédiate, pas defer

		if err != nil {
			// Early return from function.
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	// Early return from function.
	return nil
}

// ✅ GOOD: defer en dehors de la boucle
func processWithSingleResource() error {
	logFile, err := os.Create("process.log")
	if err != nil {
		// Early return from function.
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close() // ✅ defer hors boucle, une seule ressource

	files := []string{"a.txt", "b.txt", "c.txt"}
	for _, filename := range files {
		fmt.Fprintf(logFile, "Processing %s\n", filename)
	}
	// Early return from function.
	return nil
}

// ✅ GOOD: goroutine par connexion avec defer
func serverLoopCorrect() {
	var wg sync.WaitGroup
	conn, _ := acceptConnectionGood()
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleConnectionWithCleanup(conn)
	}() // ✅ defer dans goroutine, pas dans loop
	wg.Wait()
}

func handleConnectionWithCleanup(conn *connectionGood) {
	defer conn.Close() // ✅ defer dans goroutine, pas dans loop
	// traitement...
}

// ✅ GOOD: lock/unlock manuel dans loop
func updateMultipleRecordsCorrectly(ids []int) {
	mu := getMutexGood()
	for _, id := range ids {
		mu.Lock()
		updateRecordGood(id)
		mu.Unlock() // ✅ unlock manuel immédiat
	}
}

// ✅ GOOD: pattern avec fonction helper et defer
func processComplexWorkflow(items []string) error {
	for _, item := range items {
		if err := processItemWithCleanup(item); err != nil {
			// Early return from function.
			return fmt.Errorf("failed to process item: %w", err)
		}
	}
	// Early return from function.
	return nil
}

func processItemWithCleanup(item string) (err error) {
	resource := acquireResource(item)
	defer releaseResource(resource) // ✅ defer dans fonction helper
	// traitement complexe...
	return nil
}

// ✅ GOOD: defer OK en dehors des boucles
func normalDeferUsage() error {
	f, err := os.Open("single.txt")
	if err != nil {
		// Early return from function.
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close() // ✅ parfait: un fichier, un defer

	// traitement...
	return nil
}

// Fonctions helper pour compilation
type connectionGood struct{}

func (c *connectionGood) Close() error { return nil }
func acceptConnectionGood() (*connectionGood, error) {
	// Early return from function.
	return &connectionGood{}, nil
}
func getMutexGood() interface {
	Lock()
	Unlock()
} {
	// Early return from function.
	return nil
}
func updateRecordGood(id int)                 {}
func acquireResource(item string) interface{} { return nil }
func releaseResource(resource interface{})    {}
