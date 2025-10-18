package test002_edgecases

// Fichier avec GenDecl qui n'est pas ok (pas TypeSpec valide)
// Pour tester la branche: genDecl, ok := decl.(*ast.GenDecl); !ok
import "fmt"

var _ = fmt.Println
