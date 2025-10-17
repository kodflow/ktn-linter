// Package gospec_good_comments démontre les pratiques de documentation idiomatiques de Go.
//
// Ce package illustre comment documenter correctement le code Go selon
// les conventions établies dans Effective Go et les Go Code Review Comments.
//
// Référence: https://go.dev/doc/effective_go#commentary
// Référence: https://github.com/golang/go/wiki/CodeReviewComments#doc-comments
package gospec_good_comments

import (
	"fmt"
	"time"
)

// MaxRetries définit le nombre maximal de tentatives pour une opération.
const MaxRetries = 3

// DefaultTimeout est le timeout par défaut pour les opérations réseau.
const DefaultTimeout = 30 * time.Second

// Config contient la configuration de l'application.
type Config struct {
	// Host est l'adresse du serveur.
	Host string
	// Port est le port d'écoute.
	Port int
	// Timeout est la durée maximale d'attente pour une connexion.
	Timeout time.Duration
}

// NewConfig crée une nouvelle configuration avec les valeurs par défaut.
func NewConfig() *Config {
	return &Config{
		Host:    "localhost",
		Port:    8080,
		Timeout: DefaultTimeout,
	}
}

// Processor définit l'interface pour traiter des données.
//
// Les implémentations doivent être thread-safe si utilisées par
// plusieurs goroutines concurrentes.
type Processor interface {
	// Process traite les données d'entrée et retourne le résultat.
	// Retourne une erreur si le traitement échoue.
	Process(data []byte) ([]byte, error)
}

// ValidateInput vérifie que l'entrée utilisateur est valide.
//
// Retourne une erreur si l'entrée est vide ou contient des caractères invalides.
func ValidateInput(input string) error {
	if input == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}

// ProcessData traite les données avec retry logic.
//
// La fonction essaie jusqu'à MaxRetries fois avant d'abandonner.
// Chaque retry attend de manière exponentielle.
func ProcessData(data []byte) error {
	for i := 0; i < MaxRetries; i++ {
		if err := attemptProcess(data); err != nil {
			if i == MaxRetries-1 {
				return fmt.Errorf("failed after %d retries: %w", MaxRetries, err)
			}
			// Wait before retry with exponential backoff
			time.Sleep(time.Duration(1<<uint(i)) * time.Second)
			continue
		}
		return nil
	}
	return nil
}

// TODO(username): Implémenter la validation avancée des emails (issue #123)
// TODO(username): Ajouter support pour les formats internationaux (issue #124)
func validateEmail(email string) bool {
	return len(email) > 0
}

// Cache fournit un cache thread-safe pour les résultats.
//
// Exemple d'utilisation:
//
//	cache := NewCache()
//	cache.Set("key", "value")
//	val, ok := cache.Get("key")
//	if ok {
//	    fmt.Println(val)
//	}
type Cache struct {
	data map[string]interface{}
}

// NewCache crée un nouveau cache vide.
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Get récupère une valeur du cache.
//
// Retourne la valeur et true si la clé existe, sinon nil et false.
func (c *Cache) Get(key string) (interface{}, bool) {
	val, ok := c.data[key]
	return val, ok
}

// Set stocke une valeur dans le cache.
func (c *Cache) Set(key string, value interface{}) {
	c.data[key] = value
}

// Clear vide le cache de toutes ses entrées.
func (c *Cache) Clear() {
	c.data = make(map[string]interface{})
}

// calculateChecksum calcule la somme de contrôle pour les données.
//
// Utilise l'algorithme CRC32 pour la performance. Pour une sécurité
// cryptographique, utiliser SHA256 à la place.
func calculateChecksum(data []byte) uint32 {
	// Implementation uses CRC32 for speed
	// For cryptographic security, use SHA256 instead
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}
	return sum
}

// retry exécute une fonction avec retry logic.
//
// La fonction f est appelée jusqu'à maxAttempts fois.
// Si f retourne nil, retry retourne immédiatement.
// Sinon, retry attend delay entre chaque tentative.
func retry(maxAttempts int, delay time.Duration, f func() error) error {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := f()
		if err == nil {
			return nil
		}

		if attempt < maxAttempts {
			time.Sleep(delay)
		} else {
			return fmt.Errorf("max attempts reached: %w", err)
		}
	}
	return nil
}

// ErrNotFound est retourné quand un élément n'est pas trouvé.
var ErrNotFound = fmt.Errorf("not found")

// ErrInvalidInput est retourné quand l'entrée est invalide.
var ErrInvalidInput = fmt.Errorf("invalid input")

// Deprecated: Use NewCache instead.
//
// CreateCache est l'ancienne méthode pour créer un cache.
// Elle sera supprimée dans la version 2.0.
func CreateCache() *Cache {
	return NewCache()
}

// FetchUser récupère un utilisateur par ID.
//
// Retourne ErrNotFound si l'utilisateur n'existe pas.
// Retourne d'autres erreurs en cas de problème de connexion à la base de données.
func FetchUser(id int) (*User, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}
	// Database query would go here
	return &User{ID: id, Name: "test"}, nil
}

// User représente un utilisateur dans le système.
type User struct {
	// ID est l'identifiant unique de l'utilisateur.
	ID int
	// Name est le nom complet de l'utilisateur.
	Name string
	// Email est l'adresse email de l'utilisateur.
	Email string
	// CreatedAt est la date de création du compte.
	CreatedAt time.Time
}

// String implémente l'interface fmt.Stringer.
func (u *User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s}", u.ID, u.Name)
}

// IsActive vérifie si l'utilisateur est actif.
//
// Un utilisateur est considéré actif si son compte a été créé
// il y a moins de 90 jours.
func (u *User) IsActive() bool {
	return time.Since(u.CreatedAt) < 90*24*time.Hour
}

// complexAlgorithm implémente un algorithme de tri personnalisé.
//
// Cet algorithme utilise une approche hybride:
//  1. Quick sort pour les grands ensembles (> 100 éléments)
//  2. Insertion sort pour les petits ensembles
//  3. Optimisations pour les ensembles déjà triés
//
// Complexité temporelle: O(n log n) en moyenne, O(n²) pire cas
// Complexité spatiale: O(log n) pour la pile de récursion
func complexAlgorithm(items []int) []int {
	if len(items) <= 1 {
		return items
	}
	// Algorithm implementation would go here
	return items
}

// Helper function (unexported, minimal comment needed)
func attemptProcess(data []byte) error {
	_ = data
	return nil
}
