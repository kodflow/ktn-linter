// Fichier de test pour vérifier les fonctions spéciales Go.
package func014_special

import "fmt"

// Command est une structure avec un callback.
// Utilisée pour simuler des commandes CLI avec des callbacks.
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
//
// Returns:
//   - error: erreur éventuelle
func run() error {
	// Exécution
	return nil
}

var (
	// rootCmd utilise run comme callback.
	rootCmd = &Command{
		RunE: run, // ← run est utilisée comme callback
	}

	// helperFunc est assignée à une variable.
	helperFunc func() string = helper
)

// Execute exécute la commande.
//
// Returns:
//   - error: erreur éventuelle
func Execute() error {
	// Exécution du callback
	return rootCmd.RunE()
}

// helper est utilisée dans une assignation.
//
// Returns:
//   - string: message helper
func helper() string {
	// Retour du message
	return "helper"
}

// GetHelper retourne le helper.
//
// Returns:
//   - string: message
func GetHelper() string {
	// Appel du helper via la variable
	return helperFunc()
}

// Mux simule un ServeMux HTTP.
// Permet d'enregistrer des handlers HTTP sur des patterns.
type Mux struct{}

// MuxInterface définit les méthodes publiques de Mux.
type MuxInterface interface {
	HandleFunc(pattern string, handler func())
}

// NewMux crée une nouvelle instance de Mux.
//
// Returns:
//   - *Mux: nouvelle instance
func NewMux() *Mux {
	// Retour de la nouvelle instance
	return &Mux{}
}

// HandleFunc simule http.ServeMux.HandleFunc.
//
// Params:
//   - pattern: le pattern HTTP
//   - handler: le handler à enregistrer
func (m *Mux) HandleFunc(pattern string, handler func()) {
	// Simulation d'enregistrement
	_ = pattern
	_ = handler
}

// App simule une application web avec des handlers.
// Contient un multiplexeur HTTP pour router les requêtes.
type App struct {
	mux *Mux
}

// AppInterface définit les méthodes publiques de App.
type AppInterface interface {
	RegisterRoutes()
}

// NewApp crée une nouvelle instance de App.
//
// Params:
//   - mux: le multiplexeur HTTP
//
// Returns:
//   - *App: nouvelle instance
func NewApp(mux *Mux) *App {
	// Retour de la nouvelle instance
	return &App{mux: mux}
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
