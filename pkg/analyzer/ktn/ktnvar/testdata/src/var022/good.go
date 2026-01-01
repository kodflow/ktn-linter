// Package var022 contains test cases for KTN-VAR-022.
package var022

import "io"

// GoodLogger interface pour les tests.
type GoodLogger interface {
	Log(msg string)
}

// GoodService contient un champ interface directe.
// Cette structure illustre l'utilisation correcte des interfaces.
type GoodService struct {
	reader io.Reader
}

var (
	// goodHandler est une interface{} directe.
	goodHandler interface{}
	// goodAny est un any direct.
	goodAny any
)

// goodProcess utilise io.Reader directement.
//
// Params:
//   - _r: interface directe (bonne pratique)
func goodProcess(_r io.Reader) {
	// Paramètre ignoré car fonction de test
}

// goodHandle utilise io.Writer directement.
//
// Params:
//   - _w: interface directe (bonne pratique)
func goodHandle(_w io.Writer) {
	// Paramètre ignoré car fonction de test
}

// newGoodService crée un nouveau GoodService.
//
// Returns:
//   - *GoodService: nouvelle instance
func newGoodService() *GoodService {
	// Retourne une nouvelle instance
	return &GoodService{}
}

// init utilise les variables et fonctions définies.
func init() {
	// Appel de goodProcess
	goodProcess(nil)
	// Appel de goodHandle
	goodHandle(nil)
	// Utilisation de goodHandler
	_ = goodHandler
	// Utilisation de goodAny
	_ = goodAny
	// Création de GoodService
	_ = newGoodService()
}
