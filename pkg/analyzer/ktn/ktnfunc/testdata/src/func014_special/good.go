// Fichier de test pour vérifier les fonctions spéciales Go.
package func014_special

import "fmt"

// Command est une structure avec un callback.
type Command struct {
	RunE func() error
}

// main doit être ignorée (point d'entrée).
func main() {
	// Appel de Execute
	_ = Execute()
}

// init doit être ignorée (appelée automatiquement).
func init() {
	// Initialisation
	fmt.Println("init")
}

// run est passée comme callback - doit être détectée.
func run() error {
	// Exécution
	return nil
}

// rootCmd utilise run comme callback.
var rootCmd = &Command{
	RunE: run, // ← run est utilisée comme callback
}

// Execute exécute la commande.
//
// Returns:
//   - error: erreur éventuelle
func Execute() error {
	// Exécution du callback
	return rootCmd.RunE()
}

// helper est utilisée dans une assignation.
func helper() string {
	// Retour du message
	return "helper"
}

// usageHelper est assignée à une variable.
var helperFunc = helper

// GetHelper retourne le helper.
//
// Returns:
//   - string: message
func GetHelper() string {
	// Appel du helper via la variable
	return helperFunc()
}

// Mux simule un ServeMux HTTP.
type Mux struct{}

// HandleFunc simule http.ServeMux.HandleFunc.
func (m *Mux) HandleFunc(pattern string, handler func()) {
	// Simulation d'enregistrement
	_ = pattern
	_ = handler
}

// App simule une application web avec des handlers.
type App struct {
	mux *Mux
}

// handleLiveness est un handler HTTP passé comme argument.
func (a *App) handleLiveness() {
	// Handler de liveness
}

// handleReadiness est un handler HTTP passé comme argument.
func (a *App) handleReadiness() {
	// Handler de readiness
}

// registerHandlers est passé comme fonction à un appel.
func registerHandlers() {
	// Enregistrement des handlers
}

// RegisterRoutes enregistre les routes en passant les handlers comme arguments.
func (a *App) RegisterRoutes() {
	// Les handlers sont passés comme arguments à HandleFunc
	a.mux.HandleFunc("/live", a.handleLiveness)
	a.mux.HandleFunc("/ready", a.handleReadiness)
	// Fonction passée comme argument
	a.mux.HandleFunc("/register", registerHandlers)
}
