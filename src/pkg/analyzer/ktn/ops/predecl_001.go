package ktn_ops

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var RulePredecl001 = &analysis.Analyzer{
	Name: "KTN_PREDECL_002",
	Doc:  "Détecte le shadowing d'identifiants prédéclarés",
	Run:  runRulePredecl001,
}

var predeclaredIdentifiers = map[string]bool{
	// Types
	"bool": true, "byte": true, "complex64": true, "complex128": true,
	"error": true, "float32": true, "float64": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"rune": true, "string": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true, "uintptr": true,
	// Constants
	"true": true, "false": true, "iota": true, "nil": true,
	// Functions
	"append": true, "cap": true, "close": true, "complex": true, "copy": true,
	"delete": true, "imag": true, "len": true, "make": true, "new": true,
	"panic": true, "print": true, "println": true, "real": true, "recover": true,
}

func runRulePredecl001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok {
				return true
			}

			if decl.Tok != token.TYPE && decl.Tok != token.VAR && decl.Tok != token.CONST {
				return true
			}

			for _, spec := range decl.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if predeclaredIdentifiers[s.Name.Name] {
						pass.Reportf(s.Name.Pos(),
							"[KTN-PREDECL-001] Shadowing de l'identifiant prédéclaré '%s'.\n"+
								"Redéfinir un identifiant prédéclaré (type, fonction built-in, etc.) rend le code confus.\n"+
								"Cela cache l'identifiant original dans ce scope.\n"+
								"Choisissez un nom différent.\n"+
								"Exemple:\n"+
								"  // ❌ MAUVAIS\n"+
								"  var string = \"hello\"  // Cache le type string!\n"+
								"  len := 5  // Cache la fonction len()!\n"+
								"\n"+
								"  // ✅ CORRECT\n"+
								"  var message = \"hello\"\n"+
								"  length := 5",
							s.Name.Name)
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if predeclaredIdentifiers[name.Name] {
							pass.Reportf(name.Pos(),
								"[KTN-PREDECL-001] Shadowing de l'identifiant prédéclaré '%s'.\n"+
									"Redéfinir un identifiant prédéclaré (type, fonction built-in, etc.) rend le code confus.\n"+
									"Cela cache l'identifiant original dans ce scope.\n"+
									"Choisissez un nom différent.\n"+
									"Exemple:\n"+
									"  // ❌ MAUVAIS\n"+
									"  var string = \"hello\"  // Cache le type string!\n"+
									"  len := 5  // Cache la fonction len()!\n"+
									"\n"+
									"  // ✅ CORRECT\n"+
									"  var message = \"hello\"\n"+
									"  length := 5",
								name.Name)
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
