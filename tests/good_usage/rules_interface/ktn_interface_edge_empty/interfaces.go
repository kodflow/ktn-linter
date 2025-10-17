package goodempty

// Interfaces correctement typées sans utiliser interface{}.

// Processor définit l'interface de traitement avec types spécifiques.
type Processor interface {
	// Process traite les données et retourne le résultat.
	//
	// Params:
	//   - data: données à traiter
	//
	// Returns:
	//   - string: résultat du traitement
	Process(data string) string
}

// Cache définit l'interface de cache avec types spécifiques.
// Utilise des generics ou types concrets au lieu d'interface{}.
type Cache interface {
	// Get récupère une valeur du cache.
	//
	// Params:
	//   - key: clé de la valeur
	//
	// Returns:
	//   - string: valeur associée à la clé
	//   - bool: true si la clé existe
	Get(key string) (string, bool)

	// Set stocke une valeur dans le cache.
	//
	// Params:
	//   - key: clé de la valeur
	//   - value: valeur à stocker
	Set(key string, value string)
}

// Container définit un conteneur avec type spécifique.
type Container[T any] interface {
	// Add ajoute un élément au conteneur.
	//
	// Params:
	//   - item: élément à ajouter
	Add(item T)

	// GetAll retourne tous les éléments.
	//
	// Returns:
	//   - []T: tous les éléments du conteneur
	GetAll() []T
}
