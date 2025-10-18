// Package rules_goroutine_bad contient du code violant KTN-GOROUTINE-001 et KTN-GOROUTINE-002.
package rules_goroutine_bad

import (
	"fmt"
	"time"
)

// ❌ Code violant KTN-GOROUTINE-001 : goroutines dans boucles sans limitation

// BadUnlimitedGoroutines crée des goroutines illimitées.
//
// Params:
//   - requests: liste de requêtes
func BadUnlimitedGoroutines(requests []string) {
	for _, req := range requests {
		go handleRequest(req) // Viole KTN-GOROUTINE-001
	}
}

// BadForLoopGoroutines crée goroutines dans for classique.
//
// Params:
//   - n: nombre d'itérations
func BadForLoopGoroutines(n int) {
	for i := 0; i < n; i++ {
		go processTask(i) // Viole KTN-GOROUTINE-001
	}
}

// BadNestedLoopGoroutines crée goroutines dans boucles imbriquées.
//
// Params:
//   - rows: lignes
//   - cols: colonnes
func BadNestedLoopGoroutines(rows, cols int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			go processCell(i, j) // Viole KTN-GOROUTINE-001
		}
	}
}

// BadRangeMapGoroutines crée goroutines en parcourant map.
//
// Params:
//   - data: map de données
func BadRangeMapGoroutines(data map[string]int) {
	for key, value := range data {
		go processKeyValue(key, value) // Viole KTN-GOROUTINE-001
	}
}

// BadInfiniteLoopGoroutines crée goroutines dans boucle infinie.
func BadInfiniteLoopGoroutines() {
	for {
		go monitor() // Viole KTN-GOROUTINE-001
		time.Sleep(time.Second)
	}
}

// ❌ Code violant KTN-GOROUTINE-002 : goroutines sans synchronisation

// BadNoSync lance goroutine sans synchronisation.
//
// Params:
//   - data: données à traiter
func BadNoSync(data string) {
	go func() { // Viole KTN-GOROUTINE-002
		process(data)
	}()
	// Fonction termine immédiatement, goroutine peut ne jamais finir
}

// BadMultipleNoSync lance plusieurs goroutines sans sync.
//
// Params:
//   - items: éléments à traiter
func BadMultipleNoSync(items []string) {
	for _, item := range items { // Viole aussi KTN-GOROUTINE-001
		go func(s string) { // Viole KTN-GOROUTINE-002
			process(s)
		}(item)
	}
	// Retourne immédiatement
}

// BadBackgroundTask lance une tâche de fond sans sync.
func BadBackgroundTask() {
	go func() { // Viole KTN-GOROUTINE-002
		for {
			cleanup()
			time.Sleep(time.Minute)
		}
	}()
}

// BadHTTPHandlerGoroutine lance goroutine dans handler sans sync.
//
// Params:
//   - userID: ID utilisateur
func BadHTTPHandlerGoroutine(userID string) {
	go func() { // Viole KTN-GOROUTINE-002
		sendNotification(userID)
	}()
	// Répond immédiatement au client
}

// BadFunctionWithGoroutine fonction avec goroutine non synchronisée.
//
// Params:
//   - message: message à envoyer
func BadFunctionWithGoroutine(message string) {
	go sendEmail(message) // Viole KTN-GOROUTINE-002
}

// BadGoroutineInDefer lance goroutine dans defer sans sync.
func BadGoroutineInDefer() {
	defer func() {
		go cleanup() // Viole KTN-GOROUTINE-002
	}()

	doWork()
}

// BadGoroutineInCondition lance goroutine conditionnellement sans sync.
//
// Params:
//   - shouldProcess: condition de traitement
//   - data: données
func BadGoroutineInCondition(shouldProcess bool, data string) {
	if shouldProcess {
		go processData(data) // Viole KTN-GOROUTINE-002
	}
}

// BadGoroutineInSwitch lance goroutine dans switch sans sync.
//
// Params:
//   - action: action à effectuer
func BadGoroutineInSwitch(action string) {
	switch action {
	case "process":
		go performProcess() // Viole KTN-GOROUTINE-002
	case "cleanup":
		go performCleanup() // Viole KTN-GOROUTINE-002
	}
}

// BadSelectWithGoroutine lance goroutine dans select sans sync.
//
// Params:
//   - ch1: premier channel
//   - ch2: second channel
func BadSelectWithGoroutine(ch1, ch2 chan string) {
	select {
	case msg := <-ch1:
		go handleMessage(msg) // Viole KTN-GOROUTINE-002
	case msg := <-ch2:
		go handleMessage(msg) // Viole KTN-GOROUTINE-002
	}
}

// BadRecursiveGoroutine lance goroutines récursivement sans sync.
//
// Params:
//   - depth: profondeur
func BadRecursiveGoroutine(depth int) {
	if depth > 0 {
		go BadRecursiveGoroutine(depth - 1) // Viole KTN-GOROUTINE-002
	}
}

// ❌ Code violant les DEUX règles

// BadCombinedViolation viole KTN-GOROUTINE-001 ET KTN-GOROUTINE-002.
//
// Params:
//   - tasks: tâches à exécuter
func BadCombinedViolation(tasks []string) {
	for _, task := range tasks {
		go func(t string) { // Viole KTN-GOROUTINE-001 et KTN-GOROUTINE-002
			executeTask(t)
		}(task)
	}
	// Double problème : boucle + pas de sync
}

// BadMultipleLoopsNoSync plusieurs boucles sans sync.
//
// Params:
//   - batch1: premier lot
//   - batch2: second lot
func BadMultipleLoopsNoSync(batch1, batch2 []int) {
	for _, item := range batch1 {
		go processItem(item) // Viole KTN-GOROUTINE-001 et KTN-GOROUTINE-002
	}

	for _, item := range batch2 {
		go processItem(item) // Viole KTN-GOROUTINE-001 et KTN-GOROUTINE-002
	}
}

// Fonctions helpers pour les tests

// handleRequest simule le traitement d'une requête.
//
// Params:
//   - req: requête
func handleRequest(req string) {
	fmt.Println("Handling:", req)
}

// processTask simule le traitement d'une tâche.
//
// Params:
//   - taskID: ID de la tâche
func processTask(taskID int) {
	fmt.Printf("Processing task %d\n", taskID)
}

// processCell simule le traitement d'une cellule.
//
// Params:
//   - row: ligne
//   - col: colonne
func processCell(row, col int) {
	fmt.Printf("Cell [%d,%d]\n", row, col)
}

// processKeyValue simule le traitement d'une paire clé-valeur.
//
// Params:
//   - key: clé
//   - value: valeur
func processKeyValue(key string, value int) {
	fmt.Printf("%s = %d\n", key, value)
}

// monitor simule un monitoring.
func monitor() {
	fmt.Println("Monitoring...")
}

// process simule un traitement.
//
// Params:
//   - data: données
func process(data string) {
	fmt.Println("Processing:", data)
}

// cleanup simule un nettoyage.
func cleanup() {
	fmt.Println("Cleaning up...")
}

// sendNotification simule l'envoi d'une notification.
//
// Params:
//   - userID: ID utilisateur
func sendNotification(userID string) {
	fmt.Println("Notification sent to:", userID)
}

// sendEmail simule l'envoi d'un email.
//
// Params:
//   - message: message
func sendEmail(message string) {
	fmt.Println("Email sent:", message)
}

// doWork simule du travail.
func doWork() {
	fmt.Println("Working...")
}

// processData simule le traitement de données.
//
// Params:
//   - data: données
func processData(data string) {
	fmt.Println("Processing data:", data)
}

// performProcess simule une exécution.
func performProcess() {
	fmt.Println("Performing process...")
}

// performCleanup simule un nettoyage.
func performCleanup() {
	fmt.Println("Performing cleanup...")
}

// handleMessage simule le traitement d'un message.
//
// Params:
//   - msg: message
func handleMessage(msg string) {
	fmt.Println("Handling message:", msg)
}

// executeTask simule l'exécution d'une tâche.
//
// Params:
//   - task: tâche
func executeTask(task string) {
	fmt.Println("Executing:", task)
}

// processItem simule le traitement d'un élément.
//
// Params:
//   - item: élément
func processItem(item int) {
	fmt.Printf("Processing item %d\n", item)
}
