// Package var022 contains test cases for KTN-VAR-022.
package var022

import "io"

// Logger interface pour les tests.
type Logger interface {
	Log(msg string)
}

// BadService contient un champ pointeur vers interface.
// Cette structure illustre l'anti-pattern à éviter.
type BadService struct {
	reader *io.Reader // want "KTN-VAR-022"
}

var (
	// badHandler est un pointeur vers interface{}.
	badHandler *interface{} // want "KTN-VAR-022"
	// badAny est un pointeur vers any.
	badAny *any // want "KTN-VAR-022"
)

// badProcess utilise un pointeur vers io.Reader.
//
// Params:
//   - _r: pointeur vers interface (mauvaise pratique)
func badProcess(_r *io.Reader) { // want "KTN-VAR-022"
	// Paramètre ignoré car fonction de test
}

// badHandle utilise un pointeur vers io.Writer.
//
// Params:
//   - _w: pointeur vers interface (mauvaise pratique)
func badHandle(_w *io.Writer) { // want "KTN-VAR-022"
	// Paramètre ignoré car fonction de test
}

// newBadService crée un nouveau BadService.
//
// Returns:
//   - *BadService: nouvelle instance
func newBadService() *BadService {
	// Retourne une nouvelle instance
	return &BadService{}
}

// badInit utilise les variables et fonctions définies.
func badInit() {
	// Appel de badProcess
	badProcess(nil)
	// Appel de badHandle
	badHandle(nil)
	// Utilisation de badHandler
	_ = badHandler
	// Utilisation de badAny
	_ = badAny
	// Création de BadService
	_ = newBadService()
}

// init appelle badInit pour éviter KTN-FUNC-004.
func init() {
	badInit()
}
