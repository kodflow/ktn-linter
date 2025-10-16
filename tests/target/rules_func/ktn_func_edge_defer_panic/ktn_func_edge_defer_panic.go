package rules_func

import "fmt"

// Fonctions avec defer/panic/recover correctement documentées.

// OpenResource ouvre une ressource et garantit sa fermeture via defer.
//
// Params:
//   - name: nom de la ressource à ouvrir
//
// Returns:
//   - error: erreur si le nom est invalide
func OpenResource(name string) error {
	// Defer garantit la fermeture même en cas de panic
	defer closeResource(name)

	// Validation du nom de ressource
	if name == "" {
		panic("empty resource name")
	}

	// Retourne nil car la ressource est ouverte avec succès
	return nil
}

// closeResource ferme une ressource.
//
// Params:
//   - name: nom de la ressource à fermer
func closeResource(name string) {
	fmt.Println("closing", name)
}

// ProcessWithRecover traite une opération risquée avec gestion de panic.
//
// Returns:
//   - error: erreur capturée depuis un panic éventuel
func ProcessWithRecover() (err error) {
	// Recover capte les panics et les convertit en erreurs
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// Opération qui peut paniquer
	riskyOperation()
	// Retourne nil car aucune erreur n'est survenue
	return nil
}

// riskyOperation exécute une opération risquée.
func riskyOperation() {
	panic("something went wrong")
}

// MultipleDefers démontre l'ordre d'exécution LIFO des defer.
func MultipleDefers() {
	// Les defer s'exécutent dans l'ordre inverse (LIFO: Last In, First Out)
	defer fmt.Println("first defer")  // S'exécute en dernier
	defer fmt.Println("second defer") // S'exécute en deuxième
	defer fmt.Println("third defer")  // S'exécute en premier

	fmt.Println("function body")
}

// ProcessFiles traite plusieurs fichiers avec gestion appropriée des defer.
//
// Params:
//   - files: liste des noms de fichiers à traiter
//
// Returns:
//   - error: erreur si le traitement échoue
func ProcessFiles(files []string) error {
	// Approche correcte: traiter chaque fichier dans une fonction séparée
	for _, file := range files {
		if err := processFile(file); err != nil {
			// Retourne l'erreur si le traitement du fichier échoue
			return err
		}
	}
	// Retourne nil car tous les fichiers ont été traités avec succès
	return nil
}

// processFile traite un fichier.
//
// Params:
//   - name: nom du fichier à traiter
//
// Returns:
//   - error: erreur si le traitement échoue
func processFile(name string) error {
	f := openFile(name)
	// Defer dans la fonction de traitement individuel (évite l'accumulation)
	defer closeFile(f)
	// Retourne nil car le traitement est terminé avec succès
	return nil
}

// openFile ouvre un fichier.
//
// Params:
//   - name: nom du fichier à ouvrir
//
// Returns:
//   - *file: pointeur vers le fichier ouvert
func openFile(name string) *file {
	// Retourne un pointeur vers le fichier ouvert
	return &file{name: name}
}

// closeFile ferme un fichier.
//
// Params:
//   - f: pointeur vers le fichier à fermer
func closeFile(f *file) {
	fmt.Println("closing", f.name)
}

// file représente un fichier ouvert.
type file struct {
	name string
}

// HandlePanic gère les panics de manière conditionnelle selon la configuration.
//
// Params:
//   - shouldPanic: indique si un panic est attendu
//
// Returns:
//   - error: erreur si un panic inattendu se produit
func HandlePanic(shouldPanic bool) error {
	// Gestion intelligente des panics selon le contexte
	defer func() {
		if r := recover(); r != nil {
			// Panic attendu: on log et on continue
			if shouldPanic {
				fmt.Println("expected panic")
			} else {
				// Panic inattendu: on le propage
				panic(r)
			}
		}
	}()

	// Génération conditionnelle d'un panic
	if shouldPanic {
		panic("intentional panic")
	}

	// Retourne nil car aucune erreur n'est survenue
	return nil
}
