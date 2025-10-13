package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/astutil"
	"github.com/kodflow/ktn-linter/src/internal/naming"
)

// VarAnalyzer vérifie que les variables respectent les règles KTN
var VarAnalyzer = &analysis.Analyzer{
	Name: "ktnvar",
	Doc:  "Vérifie que les variables sont regroupées, documentées et typées explicitement",
	Run:  runVarAnalyzer,
}

func runVarAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Map pour suivre les commentaires avant chaque déclaration
		comments := make(map[token.Pos]*ast.CommentGroup)
		for _, cg := range file.Comments {
			if len(cg.List) > 0 {
				comments[cg.End()] = cg
			}
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}

			// KTN-VAR-001: Vérifier si c'est une variable non groupée
			if genDecl.Lparen == token.NoPos {
				// Variable déclarée sans ()
				for _, spec := range genDecl.Specs {
					valueSpec := spec.(*ast.ValueSpec)
					for _, name := range valueSpec.Names {
						pass.Reportf(name.Pos(),
							"[KTN-VAR-001] Variable '%s' déclarée individuellement. Regroupez les variables dans un bloc var ().\nExemple:\n  var (\n      %s %s = ...\n  )",
							name.Name, name.Name, astutil.GetTypeString(valueSpec))
					}
				}
				continue
			}

			// C'est un groupe de variables var ()
			// KTN-VAR-002: Vérifier le commentaire de groupe
			hasGroupComment := false
			if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
				hasGroupComment = true
			}

			if !hasGroupComment {
				pass.Reportf(genDecl.Pos(),
					"[KTN-VAR-002] Groupe de variables sans commentaire de groupe.\nAjoutez un commentaire avant le bloc var () pour décrire l'ensemble.\nExemple:\n  // Description du groupe de variables\n  var (...)")
			}

			// Vérifier chaque variable dans le groupe
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				// Vérifier si le commentaire de la variable est le même que le commentaire de groupe
				isGroupCommentOnly := hasGroupComment &&
					valueSpec.Doc != nil &&
					genDecl.Doc != nil &&
					valueSpec.Doc.Pos() == genDecl.Doc.Pos()
				checkVarSpec(pass, valueSpec, isGroupCommentOnly)
			}
		}
	}

	return nil, nil
}

func checkVarSpec(pass *analysis.Pass, spec *ast.ValueSpec, isFirstWithGroupComment bool) {
	// KTN-VAR-006: Vérifier les déclarations multiples sur une ligne
	if len(spec.Names) > 1 {
		pass.Reportf(spec.Pos(),
			"[KTN-VAR-006] Déclaration de plusieurs variables sur une ligne.\nDéclarez chaque variable sur sa propre ligne pour plus de clarté.\nExemple:\n  var (\n      %s %s = ...\n      %s %s = ...\n  )",
			spec.Names[0].Name, astutil.GetTypeString(spec),
			spec.Names[1].Name, astutil.GetTypeString(spec))
	}

	for _, name := range spec.Names {
		if name.Name == "_" {
			continue
		}

		// KTN-VAR-008: Vérifier les underscores dans le nom (sauf _)
		if strings.Contains(name.Name, "_") {
			pass.Reportf(name.Pos(),
				"[KTN-VAR-008] Variable '%s' contient un underscore.\nUtilisez la convention MixedCaps (exemple: myVariable, MaxCount).",
				name.Name)
		}

		// KTN-VAR-009: Vérifier si le nom est en ALL_CAPS (mais autoriser initialismes Go)
		if naming.IsAllCaps(name.Name) && !naming.IsValidInitialism(name.Name) {
			pass.Reportf(name.Pos(),
				"[KTN-VAR-009] Variable '%s' est en ALL_CAPS.\nUtilisez la convention MixedCaps (exemple: maxRetries au lieu de MAX_RETRIES).\nNote: Les initialismes Go valides (HTTP, URL, ID, etc.) sont autorisés.",
				name.Name)
		}

		// KTN-VAR-003: Vérifier le commentaire individuel
		hasComment := false
		if spec.Doc != nil && len(spec.Doc.List) > 0 {
			if !isFirstWithGroupComment {
				hasComment = true
			}
		} else if spec.Comment != nil && len(spec.Comment.List) > 0 {
			hasComment = true
		}

		if !hasComment {
			pass.Reportf(name.Pos(),
				"[KTN-VAR-003] Variable '%s' sans commentaire individuel.\nChaque variable doit avoir un commentaire explicatif.\nExemple:\n  // %s décrit son rôle\n  %s %s = ...",
				name.Name, name.Name, name.Name, astutil.GetTypeString(spec))
		}

		// KTN-VAR-004: Vérifier le type explicite
		if spec.Type == nil {
			pass.Reportf(name.Pos(),
				"[KTN-VAR-004] Variable '%s' sans type explicite.\nSpécifiez toujours le type : bool, string, int, []string, map[string]int, etc.\nExemple:\n  %s int = ...",
				name.Name, name.Name)
		}

		// KTN-VAR-007: Vérifier les channels sans buffer size explicite
		if spec.Type != nil {
			if chanType, ok := spec.Type.(*ast.ChanType); ok {
				if chanType.Value != nil {
					// C'est un channel, vérifier si un buffer size est spécifié
					// On vérifie si une valeur est fournie et si c'est make()
					if len(spec.Values) > 0 {
						if callExpr, ok := spec.Values[0].(*ast.CallExpr); ok {
							if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "make" {
								// make() appelé, vérifier le nombre d'arguments
								// make(chan T) = pas de buffer (potentiel problème)
								// make(chan T, size) = buffer explicite (OK)
								// MAIS: si le commentaire mentionne "unbuffered", c'est OK
								if len(callExpr.Args) < 2 {
									// Vérifier si le commentaire mentionne "unbuffered"
									commentText := ""
									if spec.Doc != nil {
										for _, c := range spec.Doc.List {
											commentText += c.Text + " "
										}
									}
									if spec.Comment != nil {
										for _, c := range spec.Comment.List {
											commentText += c.Text + " "
										}
									}

									// Si "unbuffered" n'est pas mentionné, signaler l'erreur
									if !strings.Contains(strings.ToLower(commentText), "unbuffered") {
										pass.Reportf(name.Pos(),
											"[KTN-VAR-007] Channel '%s' déclaré sans buffer size explicite.\nSpécifiez toujours la taille du buffer : make(chan T, 0) pour unbuffered, make(chan T, N) pour buffered.\nOu mentionnez 'unbuffered' dans le commentaire si intentionnel.\nExemple:\n  %s = make(%s, 0)",
											name.Name, name.Name, astutil.ExprToString(spec.Type))
									}
								}
							}
						}
					}
				}
			}
		}

		// KTN-VAR-005: Vérifier si la variable devrait être une constante
		// Règle très conservatrice : ne signaler que les cas évidents de constantes
		// mathématiques ou valeurs scientifiques immuables (Pi, E, Euler, etc.)
		if len(spec.Values) > 0 && spec.Type != nil {
			if astutil.IsConstCompatibleType(spec.Type) && astutil.IsLiteralValue(spec.Values[0]) && astutil.LooksLikeConstantName(name.Name) {
				pass.Reportf(name.Pos(),
					"[KTN-VAR-005] Variable '%s' avec valeur littérale semble être une constante immuable.\nSi la valeur ne change jamais (ex: Pi, constantes mathématiques), utilisez 'const' au lieu de 'var'.\nExemple:\n  const %s %s = ...",
					name.Name, name.Name, astutil.GetTypeString(spec))
			}
		}
	}
}
