// Good examples for the var011 test case.
package var011

import (
	"context"
	"fmt"
	"io"
	"os"
)

// goodNoShadowing utilise la réassignation correcte (pas de shadowing).
//
// Params:
//   - path: chemin du fichier
//
// Returns:
//   - error: erreur éventuelle
func goodNoShadowing(path string) error {
	file, err := os.Open(path)
	// Vérification de la condition
	if err != nil {
		// Retour de la fonction
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	// Vérification de la condition
	if err != nil {
		// Retour de la fonction
		return err
	}

	// Vérification de la condition
	if len(data) > 0 {
		err = goodValidateData(data) // OK: réassignation avec '='
		// Vérification de la condition
		if err != nil {
			// Retour de la fonction
			return err
		}
	}

	// Retour de la fonction
	return err
}

// goodShadowingErr démontre le shadowing de err qui est exemptée.
// Le shadowing de 'err' est autorisé car c'est un pattern idiomatique Go.
//
// Params:
//   - url: URL de connexion
//
// Returns:
//   - error: erreur éventuelle
func goodShadowingErr(url string) error {
	conn, err := goodDial(url)
	// Vérification de la condition
	if err != nil {
		err := fmt.Errorf("failed to connect: %w", err) // OK: 'err' est exemptée
		_ = conn
		// Retour de la fonction
		return err
	}
	// Retour de la fonction
	return nil
}

// goodShadowingOk démontre le shadowing de ok qui est exemptée.
// Le shadowing de 'ok' est autorisé pour le pattern map access/type assertion.
func goodShadowingOk() {
	m := map[string]int{"key": 42}
	v, ok := m["key"]
	// Vérification de la condition
	if ok {
		_, ok := m["other"] // OK: 'ok' est exemptée
		_ = ok
	}
	_ = v
}

// goodShadowingCtx démontre le shadowing de ctx qui est exemptée.
// Le shadowing de 'ctx' est autorisé pour la redéfinition de context.
//
// Params:
//   - ctx: contexte parent
func goodShadowingCtx(ctx context.Context) {
	// Création d'un sous-contexte
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Bloc imbriqué
	{
		ctx := context.WithValue(ctx, "key", "value") // OK: 'ctx' est exemptée
		_ = ctx
	}
}

// goodNewVariable déclare une nouvelle variable avec un nom différent.
//
// Returns:
//   - error: erreur éventuelle
func goodNewVariable() error {
	result, err := goodDoSomething()
	// Vérification de la condition
	if err != nil {
		// Retour de la fonction
		return err
	}

	// Vérification de la condition
	if result > 0 {
		err2 := goodDoAnotherThing() // OK: nouvelle variable avec nom différent
		// Vérification de la condition
		if err2 != nil {
			// Retour de la fonction
			return err2
		}
	}

	err = goodFinalCheck() // OK: réassignation
	// Retour de la fonction
	return err
}

// goodLocalScopeErr déclare err dans un scope différent (OK - pas de parent).
//
// Returns:
//   - error: erreur éventuelle
func goodLocalScopeErr() error {
	// Vérification de la condition
	if true {
		err := goodDoSomething2() // OK: première déclaration dans ce scope
		// Vérification de la condition
		if err != nil {
			// Retour de la fonction
			return err
		}
	}

	// Vérification de la condition
	if false {
		err := goodDoAnotherThing() // OK: première déclaration dans ce scope (différent du précédent)
		// Vérification de la condition
		if err != nil {
			// Retour de la fonction
			return err
		}
	}

	// Retour de la fonction
	return nil
}

// goodValidateData valide les données.
//
// Params:
//   - _data: données à valider (non utilisées dans cet exemple)
//
// Returns:
//   - error: erreur éventuelle
func goodValidateData(_data []byte) error {
	// Retour de la fonction
	return nil
}

// goodDial établit une connexion.
//
// Params:
//   - _url: URL de connexion (non utilisée dans cet exemple)
//
// Returns:
//   - any: connexion
//   - error: erreur éventuelle
func goodDial(_url string) (any, error) {
	// Retour de la fonction
	return nil, nil
}

// goodDoSomething effectue une opération.
//
// Returns:
//   - int: résultat
//   - error: erreur éventuelle
func goodDoSomething() (int, error) {
	// Retour de la fonction
	return 0, nil
}

// goodDoSomething2 effectue une autre opération.
//
// Returns:
//   - error: erreur éventuelle
func goodDoSomething2() error {
	// Retour de la fonction
	return nil
}

// goodDoAnotherThing effectue encore une autre opération.
//
// Returns:
//   - error: erreur éventuelle
func goodDoAnotherThing() error {
	// Retour de la fonction
	return nil
}

// goodFinalCheck effectue une vérification finale.
//
// Returns:
//   - error: erreur éventuelle
func goodFinalCheck() error {
	// Retour de la fonction
	return nil
}

// init utilise les fonctions privées
func init() {
	// Appel de goodNoShadowing
	_ = goodNoShadowing("")
	// Appel de goodShadowingErr
	_ = goodShadowingErr("")
	// Appel de goodShadowingOk
	goodShadowingOk()
	// Appel de goodShadowingCtx
	goodShadowingCtx(context.Background())
	// Appel de goodNewVariable
	goodNewVariable()
	// Appel de goodLocalScopeErr
	goodLocalScopeErr()
	// Appel de goodValidateData
	_ = goodValidateData(nil)
	// Appel de goodDial
	_, _ = goodDial("")
	// Appel de goodDoSomething
	goodDoSomething()
	// Appel de goodDoSomething2
	goodDoSomething2()
	// Appel de goodDoAnotherThing
	goodDoAnotherThing()
	// Appel de goodFinalCheck
	goodFinalCheck()
}
