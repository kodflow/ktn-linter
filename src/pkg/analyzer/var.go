package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
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
							name.Name, name.Name, getTypeString(valueSpec))
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
			spec.Names[0].Name, getTypeString(spec),
			spec.Names[1].Name, getTypeString(spec))
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
		if isAllCaps(name.Name) && !isValidInitialism(name.Name) {
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
				name.Name, name.Name, name.Name, getTypeString(spec))
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
											name.Name, name.Name, exprToString(spec.Type))
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
			if isConstCompatibleType(spec.Type) && isLiteralValue(spec.Values[0]) && looksLikeConstantName(name.Name) {
				pass.Reportf(name.Pos(),
					"[KTN-VAR-005] Variable '%s' avec valeur littérale semble être une constante immuable.\nSi la valeur ne change jamais (ex: Pi, constantes mathématiques), utilisez 'const' au lieu de 'var'.\nExemple:\n  const %s %s = ...",
					name.Name, name.Name, getTypeString(spec))
			}
		}
	}
}

// isAllCaps vérifie si une chaîne est entièrement en majuscules
func isAllCaps(s string) bool {
	if len(s) == 0 {
		return false
	}
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
	return hasLetter
}

// isValidInitialism vérifie si le nom est composé uniquement d'initialismes Go valides
// Exemples valides: HTTPOK, URLID, APIURL, HTTPSPort
// Exemples invalides: MAX_BUFFER, HTTP_OK
func isValidInitialism(name string) bool {
	// Liste des initialismes Go courants (voir Effective Go)
	initialisisms := []string{
		"HTTP", "HTTPS", "URL", "URI", "ID", "API", "JSON", "XML", "HTML",
		"SQL", "TLS", "SSL", "TCP", "UDP", "IP", "DNS", "SSH", "FTP",
		"OK", "EOF", "UID", "UUID", "ASCII", "UTF", "CPU", "RAM", "IO",
		"DB", "RPC", "CDN", "AWS", "GCP", "TTL", "ACL", "CORS", "CSRF",
	}

	// Si le nom contient un underscore, c'est invalide (KTN-VAR-008)
	if strings.Contains(name, "_") {
		return false
	}

	// Essayer de décomposer le nom en initialismes connus
	remaining := name
	matched := false

	for len(remaining) > 0 {
		foundMatch := false
		// Essayer de matcher le début avec un initialisme
		for _, init := range initialisisms {
			if strings.HasPrefix(remaining, init) {
				remaining = remaining[len(init):]
				foundMatch = true
				matched = true
				break
			}
		}

		// Si on n'a pas trouvé de match et qu'il reste des caractères
		if !foundMatch {
			// Vérifier si le reste est en MixedCaps (ex: HTTPServer, URLParser)
			if remaining != "" && unicode.IsUpper(rune(remaining[0])) {
				// C'est peut-être une combinaison initialism+nom (HTTPOK, HTTPNotFound)
				// On accepte si au moins un initialisme a été trouvé
				return matched
			}
			return false
		}
	}

	return matched
}

// isConstCompatibleType vérifie si le type est compatible avec const
func isConstCompatibleType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Types de base compatibles avec const
		switch e.Name {
		case "bool", "string",
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64",
			"complex64", "complex128",
			"byte", "rune":
			return true
		}
	}
	return false
}

// isLiteralValue vérifie si l'expression est une valeur littérale
func isLiteralValue(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.BasicLit: // true, false, 123, "string", 3.14, etc.
		return true
	case *ast.Ident: // true, false, nil
		return true
	}
	return false
}

// looksLikeConstantName vérifie si le nom ressemble à une constante mathématique ou scientifique
// Ne signale que les cas évidents (Pi, E, Euler, etc.) pour éviter les faux positifs
func looksLikeConstantName(name string) bool {
	// Liste de noms connus de constantes mathématiques/scientifiques
	knownConstants := map[string]bool{
		"Pi":             true,
		"E":              true,
		"Euler":          true,
		"EulerNumber":    true,
		"GoldenRatio":    true,
		"Phi":            true,
		"Tau":            true,
		"SpeedOfLight":   true,
		"PlanckConstant": true,
		"AvogadroNumber": true,
		"BoltzmannConst": true,
		"GravityConst":   true,
	}

	// Vérifier si c'est un nom connu
	if knownConstants[name] {
		return true
	}

	// Vérifier si le nom contient des indicateurs de constante mathématique
	nameLower := strings.ToLower(name)
	if strings.Contains(nameLower, "constant") ||
		strings.Contains(nameLower, "ratio") && strings.Contains(nameLower, "golden") {
		return true
	}

	return false
}
