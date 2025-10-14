package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/astutil"
	"github.com/kodflow/ktn-linter/src/internal/naming"
)

// Analyzers
var (
	// VarAnalyzer vérifie que les variables respectent les règles KTN
	VarAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnvar",
		Doc:  "Vérifie que les variables sont regroupées, documentées et typées explicitement",
		Run:  runVarAnalyzer,
	}
)

// runVarAnalyzer vérifie que toutes les variables respectent les règles KTN.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//
// Returns:
//   - interface{}: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runVarAnalyzer(pass *analysis.Pass) (interface{}, error) {
	// Phase 1 : Collecter toutes les variables package-level et leurs assignations
	packageVars := make(map[*ast.Ident]*ast.ValueSpec)
	reassignedVars := make(map[string]bool)

	for _, file := range pass.Files {
		collectPackageVars(file, packageVars)
		collectReassignments(file, reassignedVars)
	}

	// Phase 2 : Vérifier les déclarations avec l'info des réassignations
	for _, file := range pass.Files {
		checkVarDeclarations(pass, file, reassignedVars)
	}
	return nil, nil
}

// collectPackageVars collecte toutes les variables déclarées au niveau package.
//
// Params:
//   - file: le fichier AST à parcourir
//   - packageVars: la map où stocker les variables trouvées
func collectPackageVars(file *ast.File, packageVars map[*ast.Ident]*ast.ValueSpec) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for _, name := range valueSpec.Names {
				packageVars[name] = valueSpec
			}
		}
	}
}

// collectReassignments parcourt le fichier pour trouver toutes les réassignations.
//
// Params:
//   - file: le fichier AST à parcourir
//   - reassignedVars: la map où marquer les variables réassignées
func collectReassignments(file *ast.File, reassignedVars map[string]bool) {
	ast.Inspect(file, func(n ast.Node) bool {
		// Chercher les assignations (=, pas :=)
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok || assignStmt.Tok != token.ASSIGN {
			return true
		}

		// Marquer toutes les variables du côté gauche comme réassignées
		for _, lhs := range assignStmt.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				reassignedVars[ident.Name] = true
			}
		}

		return true
	})
}

// checkVarDeclarations parcourt et vérifie toutes les déclarations de variables.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - file: le fichier AST à analyser
//   - reassignedVars: la map des variables réassignées
func checkVarDeclarations(pass *analysis.Pass, file *ast.File, reassignedVars map[string]bool) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		if genDecl.Lparen == token.NoPos {
			reportUngroupedVar(pass, genDecl)
			continue
		}

		checkVarGroup(pass, genDecl, reassignedVars)
	}
}

// reportUngroupedVar signale une variable non groupée.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - genDecl: la déclaration générale contenant la variable non groupée
func reportUngroupedVar(pass *analysis.Pass, genDecl *ast.GenDecl) {
	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		for _, name := range valueSpec.Names {
			pass.Reportf(name.Pos(),
				"[KTN-VAR-001] Variable '%s' déclarée individuellement. Regroupez les variables dans un bloc var ().\nExemple:\n  var (\n      %s %s = ...\n  )",
				name.Name, name.Name, astutil.GetTypeString(valueSpec))
		}
	}
}

// checkVarGroup vérifie un groupe de variables.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - genDecl: la déclaration de groupe var() à vérifier
//   - reassignedVars: la map des variables réassignées
func checkVarGroup(pass *analysis.Pass, genDecl *ast.GenDecl, reassignedVars map[string]bool) {
	hasGroupComment := false
	if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
		hasGroupComment = true
	}

	if !hasGroupComment {
		pass.Reportf(genDecl.Pos(),
			"[KTN-VAR-002] Groupe de variables sans commentaire de groupe.\nAjoutez un commentaire avant le bloc var () pour décrire l'ensemble.\nExemple:\n  // Description du groupe de variables\n  var (...)")
	}

	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		isGroupCommentOnly := hasGroupComment &&
			valueSpec.Doc != nil &&
			genDecl.Doc != nil &&
			valueSpec.Doc.Pos() == genDecl.Doc.Pos()
		checkVarSpec(pass, valueSpec, isGroupCommentOnly, reassignedVars)
	}
}

// checkVarSpec vérifie une spécification de variable individuelle.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - spec: la spécification de variable à vérifier
//   - isFirstWithGroupComment: true si la variable partage le commentaire de groupe
//   - reassignedVars: la map des variables réassignées
func checkVarSpec(pass *analysis.Pass, spec *ast.ValueSpec, isFirstWithGroupComment bool, reassignedVars map[string]bool) {
	checkMultipleDeclarations(pass, spec)

	for _, name := range spec.Names {
		if name.Name == "_" {
			continue
		}

		checkVarNaming(pass, name)
		checkVarComment(pass, name, spec, isFirstWithGroupComment)
		checkVarType(pass, name, spec)
		checkChannelBuffer(pass, name, spec)
		checkIfShouldBeConst(pass, name, spec, reassignedVars)
	}
}

// checkMultipleDeclarations vérifie les déclarations multiples sur une ligne.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - spec: la spécification de variable à vérifier
func checkMultipleDeclarations(pass *analysis.Pass, spec *ast.ValueSpec) {
	if len(spec.Names) > 1 {
		pass.Reportf(spec.Pos(),
			"[KTN-VAR-006] Déclaration de plusieurs variables sur une ligne.\nDéclarez chaque variable sur sa propre ligne pour plus de clarté.\nExemple:\n  var (\n      %s %s = ...\n      %s %s = ...\n  )",
			spec.Names[0].Name, astutil.GetTypeString(spec),
			spec.Names[1].Name, astutil.GetTypeString(spec))
	}
}

// checkVarNaming vérifie le nommage de la variable.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - name: l'identifiant de la variable à vérifier
func checkVarNaming(pass *analysis.Pass, name *ast.Ident) {
	if strings.Contains(name.Name, "_") {
		pass.Reportf(name.Pos(),
			"[KTN-VAR-008] Variable '%s' contient un underscore.\nUtilisez la convention MixedCaps (exemple: myVariable, MaxCount).",
			name.Name)
	}

	if naming.IsAllCaps(name.Name) && !naming.IsValidInitialism(name.Name) {
		pass.Reportf(name.Pos(),
			"[KTN-VAR-009] Variable '%s' est en ALL_CAPS.\nUtilisez la convention MixedCaps (exemple: maxRetries au lieu de MAX_RETRIES).\nNote: Les initialismes Go valides (HTTP, URL, ID, etc.) sont autorisés.",
			name.Name)
	}
}

// checkVarComment vérifie le commentaire de la variable.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - name: l'identifiant de la variable
//   - spec: la spécification de variable
//   - isFirstWithGroupComment: true si la variable partage le commentaire de groupe
func checkVarComment(pass *analysis.Pass, name *ast.Ident, spec *ast.ValueSpec, isFirstWithGroupComment bool) {
	if !hasIndividualVarComment(spec, isFirstWithGroupComment) {
		pass.Reportf(name.Pos(),
			"[KTN-VAR-003] Variable '%s' sans commentaire individuel.\nChaque variable doit avoir un commentaire explicatif.\nExemple:\n  // %s décrit son rôle\n  %s %s = ...",
			name.Name, name.Name, name.Name, astutil.GetTypeString(spec))
	}
}

// hasIndividualVarComment vérifie si une variable a un commentaire individuel.
//
// Params:
//   - spec: la spécification de variable à vérifier
//   - isFirstWithGroupComment: true si la variable partage le commentaire de groupe
//
// Returns:
//   - bool: true si un commentaire individuel existe
func hasIndividualVarComment(spec *ast.ValueSpec, isFirstWithGroupComment bool) bool {
	if spec.Doc != nil && len(spec.Doc.List) > 0 {
		if !isFirstWithGroupComment {
			return true
		}
	} else if spec.Comment != nil && len(spec.Comment.List) > 0 {
		return true
	}
	return false
}

// checkVarType vérifie le type explicite de la variable.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - name: l'identifiant de la variable
//   - spec: la spécification de variable
func checkVarType(pass *analysis.Pass, name *ast.Ident, spec *ast.ValueSpec) {
	if spec.Type == nil {
		pass.Reportf(name.Pos(),
			"[KTN-VAR-004] Variable '%s' sans type explicite.\nSpécifiez toujours le type : bool, string, int, []string, map[string]int, etc.\nExemple:\n  %s int = ...",
			name.Name, name.Name)
	}
}

// checkChannelBuffer vérifie les channels sans buffer size explicite.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - name: l'identifiant de la variable channel
//   - spec: la spécification de variable
func checkChannelBuffer(pass *analysis.Pass, name *ast.Ident, spec *ast.ValueSpec) {
	if !isMakeChannelWithoutBuffer(spec) {
		return
	}

	if isUnbufferedIntentional(spec) {
		return
	}

	pass.Reportf(name.Pos(),
		"[KTN-VAR-007] Channel '%s' déclaré sans buffer size explicite.\nSpécifiez toujours la taille du buffer : make(chan T, 0) pour unbuffered, make(chan T, N) pour buffered.\nOu mentionnez 'unbuffered' dans le commentaire si intentionnel.\nExemple:\n  %s = make(%s, 0)",
		name.Name, name.Name, astutil.ExprToString(spec.Type))
}

// isMakeChannelWithoutBuffer vérifie si c'est un channel déclaré avec make sans buffer.
//
// Params:
//   - spec: la spécification de variable à vérifier
//
// Returns:
//   - bool: true si c'est un channel make sans buffer
func isMakeChannelWithoutBuffer(spec *ast.ValueSpec) bool {
	if spec.Type == nil {
		return false
	}

	chanType, ok := spec.Type.(*ast.ChanType)
	if !ok || chanType.Value == nil {
		return false
	}

	if len(spec.Values) == 0 {
		return false
	}

	callExpr, ok := spec.Values[0].(*ast.CallExpr)
	if !ok {
		return false
	}

	ident, ok := callExpr.Fun.(*ast.Ident)
	return ok && ident.Name == "make" && len(callExpr.Args) < 2
}

// isUnbufferedIntentional vérifie si unbuffered est mentionné dans le commentaire.
//
// Params:
//   - spec: la spécification de variable avec potentiellement un commentaire
//
// Returns:
//   - bool: true si le mot "unbuffered" est présent dans le commentaire
func isUnbufferedIntentional(spec *ast.ValueSpec) bool {
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
	return strings.Contains(strings.ToLower(commentText), "unbuffered")
}

// checkIfShouldBeConst vérifie si la variable devrait être une constante.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - name: l'identifiant de la variable
//   - spec: la spécification de variable
//   - reassignedVars: la map des variables réassignées
func checkIfShouldBeConst(pass *analysis.Pass, name *ast.Ident, spec *ast.ValueSpec, reassignedVars map[string]bool) {
	// Ignorer si pas de valeur initiale ou pas de type
	if len(spec.Values) == 0 || spec.Type == nil {
		return
	}

	// Ignorer si le type n'est pas compatible avec const
	if !astutil.IsConstCompatibleType(spec.Type) {
		return
	}

	// Ignorer si la valeur n'est pas littérale
	if !astutil.IsLiteralValue(spec.Values[0]) {
		return
	}

	// Vérifier si la variable est jamais réassignée
	if reassignedVars[name.Name] {
		return
	}

	// Si la variable n'est jamais réassignée et a une valeur littérale, suggérer const
	pass.Reportf(name.Pos(),
		"[KTN-VAR-005] Variable '%s' jamais réassignée avec valeur littérale.\nCette variable se comporte comme une constante. Utilisez 'const' au lieu de 'var' pour indiquer l'intention d'immuabilité.\nExemple:\n  const %s %s = ...",
		name.Name, name.Name, astutil.GetTypeString(spec))
}
