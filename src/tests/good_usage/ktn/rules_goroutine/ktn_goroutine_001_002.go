// Package rules_goroutine_good contient du code conforme à KTN-GOROUTINE-001 et KTN-GOROUTINE-002.
package rules_goroutine_good

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ✅ Code conforme KTN-GOROUTINE-001 : goroutines contrôlées avec worker pool

// GoodWorkerPool utilise un worker pool pour limiter les goroutines.
//
// Params:
//   - requests: liste de requêtes
func GoodWorkerPool(requests []string) {
	// Worker pool avec 10 workers maximum
	const numWorkers = 10
	jobs := make(chan string, 100)
	var wg sync.WaitGroup

	// Lancer les workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Envoyer les jobs
	for _, req := range requests {
		jobs <- req
	}
	closeJobs(jobs)

	// Attendre la fin
	wg.Wait()
}

// GoodBufferedChannel utilise un buffered channel pour contrôler concurrence.
//
// Params:
//   - tasks: tâches à exécuter
func GoodBufferedChannel(tasks []int) {
	const maxConcurrency = 5
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		sem <- struct{}{} // Acquérir le sémaphore
		go semaphoreWorker(task, sem, &wg)
	}

	// Attendre toutes les goroutines
	wg.Wait()
}

// GoodContextTimeout utilise context pour limiter durée.
//
// Params:
//   - requests: requêtes
func GoodContextTimeout(requests []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	results := make(chan string, len(requests))

	for _, req := range requests {
		wg.Add(1)
		go contextWorker(ctx, req, results, &wg)
	}

	// Attendre la fin
	wg.Wait()
	closeResults(results)
}

// ✅ Code conforme KTN-GOROUTINE-002 : goroutines avec synchronisation

// GoodWithWaitGroup utilise sync.WaitGroup pour synchroniser.
//
// Params:
//   - data: données à traiter
func GoodWithWaitGroup(data string) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		process(data)
	}()

	// Attend que la goroutine termine
	wg.Wait()
}

// GoodWithContext utilise context.Context pour synchronisation.
//
// Params:
//   - ctx: contexte
//   - items: éléments à traiter
func GoodWithContext(ctx context.Context, items []string) {
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go contextProcessor(ctx, item, &wg)
	}

	// Attendre la fin
	wg.Wait()
}

// GoodWithChannel utilise channel pour synchronisation.
//
// Params:
//   - tasks: tâches à exécuter
//
// Returns:
//   - []string: résultats
func GoodWithChannel(tasks []string) []string {
	results := make(chan string, len(tasks))
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go channelWorker(task, results, &wg)
	}

	// Attendre et fermer channel
	closeResultsWhenDone(&wg, results)

	// Collecter résultats
	var output []string
	for result := range results {
		output = append(output, result)
	}

	// Retourne les résultats
	return output
}

// GoodBackgroundWithContext lance tâche de fond avec context.
//
// Params:
//   - ctx: contexte
func GoodBackgroundWithContext(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Contexte annulé, arrêter
				return
			case <-ticker.C:
				cleanup()
			}
		}
	}()
}

// GoodHTTPHandlerWithWaitGroup gère goroutine dans handler avec sync.
//
// Params:
//   - userID: ID utilisateur
func GoodHTTPHandlerWithWaitGroup(userID string) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendNotification(userID)
	}()

	// Attendre avant de répondre au client
	wg.Wait()
}

// GoodConditionalGoroutine lance goroutine conditionnellement avec sync.
//
// Params:
//   - shouldProcess: condition
//   - data: données
func GoodConditionalGoroutine(shouldProcess bool, data string) {
	if shouldProcess {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			processData(data)
		}()
		wg.Wait()
	}
}

// GoodPipelinePattern utilise pattern pipeline avec channels.
//
// Params:
//   - input: données d'entrée
//
// Returns:
//   - <-chan string: channel de résultats
func GoodPipelinePattern(input []int) <-chan string {
	// Stage 1: Convertir en string
	stage1 := pipelineStage1(input)

	// Stage 2: Transformer
	stage2 := pipelineStage2(stage1)

	// Retourne le channel de sortie
	return stage2
}

// pipelineStage1 convertit les nombres en strings.
//
// Params:
//   - input: données d'entrée
//
// Returns:
//   - <-chan string: channel de sortie
func pipelineStage1(input []int) <-chan string {
	out := make(chan string, len(input))
	go func() {
		defer close(out)
		for _, num := range input {
			out <- fmt.Sprintf("%d", num)
		}
	}()
	// Early return from function.
	return out
}

// pipelineStage2 transforme les strings.
//
// Params:
//   - in: channel d'entrée
//
// Returns:
//   - <-chan string: channel de sortie
func pipelineStage2(in <-chan string) <-chan string {
	out := make(chan string, 100)
	go func() {
		defer close(out)
		for s := range in {
			out <- "processed-" + s
		}
	}()
	// Early return from function.
	return out
}

// GoodFanOutFanIn implémente pattern fan-out/fan-in.
//
// Params:
//   - tasks: tâches à distribuer
//
// Returns:
//   - []string: résultats
func GoodFanOutFanIn(tasks []string) []string {
	const numWorkers = 3
	jobs := make(chan string, len(tasks))
	results := make(chan string, len(tasks))
	var wg sync.WaitGroup

	startWorkers(numWorkers, jobs, results, &wg)
	sendJobs(tasks, jobs)
	closeResultsWhenDone(&wg, results)

	// Retourne les résultats collectés
	return collectResults(results)
}

// startWorkers lance les workers du fan-out.
//
// Params:
//   - numWorkers: nombre de workers
//   - jobs: channel de jobs
//   - results: channel de résultats
//   - wg: waitgroup pour synchronisation
func startWorkers(numWorkers int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go fanOutTaskWorker(jobs, results, wg)
	}
}

// fanOutTaskWorker traite les tâches du channel jobs.
//
// Params:
//   - jobs: channel de jobs
//   - results: channel de résultats
//   - wg: waitgroup
func fanOutTaskWorker(jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range jobs {
		results <- processTask2(task)
	}
}

// sendJobs envoie les jobs sur le channel.
//
// Params:
//   - tasks: tâches à envoyer
//   - jobs: channel de jobs
func sendJobs(tasks []string, jobs chan<- string) {
	go func() {
		for _, task := range tasks {
			jobs <- task
		}
		close(jobs)
	}()
}

// closeResultsWhenDone ferme le channel results quand tous les workers terminent.
//
// Params:
//   - wg: waitgroup des workers
//   - results: channel à fermer
func closeResultsWhenDone(wg *sync.WaitGroup, results chan string) {
	go func() {
		wg.Wait()
		close(results)
	}()
}

// collectResults collecte les résultats depuis un channel.
//
// Params:
//   - results: channel de résultats
//
// Returns:
//   - []string: résultats collectés
func collectResults(results <-chan string) []string {
	var output []string
	for result := range results {
		output = append(output, result)
	}
	// Retourne les résultats collectés
	return output
}

// GoodGracefulShutdown implémente arrêt gracieux avec context.
//
// Params:
//   - ctx: contexte principal
func GoodGracefulShutdown(ctx context.Context) {
	var wg sync.WaitGroup

	// Lancer plusieurs workers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go shutdownWorker(ctx, i, &wg)
	}

	// Attendre que tous les workers terminent
	wg.Wait()
}

// GoodErrorHandling gère les erreurs dans goroutines.
//
// Params:
//   - tasks: tâches
//
// Returns:
//   - error: première erreur rencontrée
func GoodErrorHandling(tasks []string) error {
	errCh := make(chan error, len(tasks))
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go errorWorker(task, errCh, &wg)
	}

	// Attendre et fermer
	closeErrorChannelWhenDone(&wg, errCh)

	// Vérifier les erreurs
	for err := range errCh {
		if err != nil {
			// Retourne la première erreur wrappée avec contexte
			return fmt.Errorf("goroutine execution failed: %w", err)
		}
	}

	// Retourne nil si succès
	return nil
}

// Fonctions helpers pour les tests

// semaphoreWorker traite une tâche avec sémaphore.
//
// Params:
//   - t: ID de la tâche
//   - sem: channel sémaphore
//   - wg: waitgroup
func semaphoreWorker(t int, sem chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		<-sem // Libérer le sémaphore
		wg.Done()
	}()
	processTask(t)
}

// contextWorker traite une requête avec contexte et channel.
//
// Params:
//   - ctx: contexte
//   - r: requête
//   - results: channel de résultats
//   - wg: waitgroup
func contextWorker(ctx context.Context, r string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		// Contexte annulé
		return
	case results <- processRequest(r):
		// Traitement réussi
	}
}

// closeResults ferme le channel de résultats.
//
// Params:
//   - results: channel à fermer
func closeResults(results chan string) {
	close(results)
}

// contextProcessor traite un élément avec contexte.
//
// Params:
//   - ctx: contexte
//   - s: élément à traiter
//   - wg: waitgroup
func contextProcessor(ctx context.Context, s string, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		// Contexte annulé, arrêter
		return
	default:
		process(s)
	}
}

// channelWorker traite une tâche et envoie le résultat au channel.
//
// Params:
//   - t: tâche
//   - results: channel de résultats
//   - wg: waitgroup
func channelWorker(t string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	result := processTask2(t)
	results <- result
}

// shutdownWorker implémente un worker avec arrêt gracieux.
//
// Params:
//   - ctx: contexte
//   - workerID: ID du worker
//   - wg: waitgroup
func shutdownWorker(ctx context.Context, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// Signal d'arrêt reçu
			fmt.Printf("Worker %d shutting down\n", workerID)
			// Retourne pour terminer proprement
			return
		default:
			// Continue le travail
			doWork()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// errorWorker effectue une opération risquée et envoie les erreurs au channel.
//
// Params:
//   - t: tâche
//   - errCh: channel d'erreurs
//   - wg: waitgroup
func errorWorker(t string, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := riskyOperation(t); err != nil {
		errCh <- err
	}
}

// closeErrorChannelWhenDone ferme le channel d'erreurs quand tous les workers terminent.
//
// Params:
//   - wg: waitgroup des workers
//   - errCh: channel à fermer
func closeErrorChannelWhenDone(wg *sync.WaitGroup, errCh chan error) {
	go func() {
		wg.Wait()
		close(errCh)
	}()
}

// worker traite les jobs depuis un channel.
//
// Params:
//   - id: ID du worker
//   - jobs: channel de jobs
//   - wg: waitgroup
func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range jobs {
		handleRequest(req)
	}
}

// closeJobs ferme le channel de jobs.
//
// Params:
//   - jobs: channel à fermer
func closeJobs(jobs chan string) {
	close(jobs)
}

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

// processRequest simule le traitement d'une requête.
//
// Params:
//   - req: requête
//
// Returns:
//   - string: résultat
func processRequest(req string) string {
	// Retourne le résultat
	return "processed: " + req
}

// process simule un traitement.
//
// Params:
//   - data: données
func process(data string) {
	fmt.Println("Processing:", data)
}

// processTask2 simule le traitement d'une tâche.
//
// Params:
//   - task: tâche
//
// Returns:
//   - string: résultat
func processTask2(task string) string {
	// Retourne le résultat
	return "result-" + task
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

// processData simule le traitement de données.
//
// Params:
//   - data: données
func processData(data string) {
	fmt.Println("Processing data:", data)
}

// doWork simule du travail.
func doWork() {
	fmt.Println("Working...")
}

// riskyOperation simule une opération risquée.
//
// Params:
//   - task: tâche
//
// Returns:
//   - error: erreur potentielle
func riskyOperation(task string) error {
	// Retourne nil car succès
	return nil
}
