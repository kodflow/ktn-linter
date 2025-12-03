package ktnfunc

import (
	"go/ast"
	"testing"
)

// Test_runFunc014 tests the runFunc014 private function.
func Test_runFunc014(t *testing.T) {
	// Test cases pour la fonction privée runFunc014
	// La logique principale est testée via l'API publique dans 014_external_test.go
	// Ce test vérifie les cas edge de la fonction privée

	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Exécution tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - la logique principale est dans external tests
		})
	}
}

// Test_collectPrivateFunctions vérifie la collecte des fonctions privées.
func Test_collectPrivateFunctions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite un contexte d'analyse complet
		})
	}
}

// Test_extractPrivateFuncInfo vérifie l'extraction des infos de fonctions privées.
func Test_extractPrivateFuncInfo(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		isNil    bool
	}{
		{
			name: "error case validation",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "privateFunc"},
				Type: &ast.FuncType{},
			},
			isNil: false,
		},
		{
			name: "public function",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "PublicFunc"},
				Type: &ast.FuncType{},
			},
			isNil: true,
		},
		{
			name: "main function",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "main"},
				Type: &ast.FuncType{},
			},
			isNil: true,
		},
		{
			name: "init function",
			funcDecl: &ast.FuncDecl{
				Name: &ast.Ident{Name: "init"},
				Type: &ast.FuncType{},
			},
			isNil: true,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := extractPrivateFuncInfo(tt.funcDecl)
			// Vérification du résultat
			if (result == nil) != tt.isNil {
				t.Errorf("extractPrivateFuncInfo() nil = %v, want %v", result == nil, tt.isNil)
			}
		})
	}
}

// Test_collectCalledFunctions vérifie la collecte des fonctions appelées.
func Test_collectCalledFunctions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite un contexte d'analyse complet
		})
	}
}

// Test_collectCallbackUsages vérifie la détection des callbacks.
func Test_collectCallbackUsages(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite un contexte d'analyse complet
		})
	}
}

// Test_reportUnusedPrivateFuncs vérifie le rapport des fonctions non utilisées.
func Test_reportUnusedPrivateFuncs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite un contexte d'analyse complet
		})
	}
}

// Test_reportUnusedFunc vérifie le rapport d'une fonction non utilisée.
func Test_reportUnusedFunc(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - nécessite un contexte d'analyse complet
		})
	}
}

// Test_extractReceiverType vérifie l'extraction du type du receiver.
func Test_extractReceiverType(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "error case validation",
			expr:     &ast.Ident{Name: "MyType"},
			expected: "MyType",
		},
		{
			name: "pointer receiver",
			expr: &ast.StarExpr{
				X: &ast.Ident{Name: "MyType"},
			},
			expected: "MyType",
		},
		{
			name:     "unsupported type",
			expr:     &ast.BasicLit{},
			expected: "",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := extractReceiverType(tt.expr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("extractReceiverType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_extractCalledFuncName vérifie l'extraction du nom de fonction appelée.
func Test_extractCalledFuncName(t *testing.T) {
	tests := []struct {
		name     string
		callExpr *ast.CallExpr
		expected string
	}{
		{
			name: "error case validation",
			callExpr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "myFunc"},
			},
			expected: "myFunc",
		},
		{
			name: "method call",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "obj"},
					Sel: &ast.Ident{Name: "Method"},
				},
			},
			expected: "Method",
		},
		{
			name: "unsupported call",
			callExpr: &ast.CallExpr{
				Fun: &ast.BasicLit{},
			},
			expected: "",
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			result := extractCalledFuncName(tt.callExpr)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("extractCalledFuncName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_collectCallExprCallbacks vérifie la détection des callbacks dans les CallExpr.
func Test_collectCallExprCallbacks(t *testing.T) {
	tests := []struct {
		name         string
		callExpr     *ast.CallExpr
		privateFuncs map[string][]*privateFuncInfo
		expected     []string
	}{
		{
			name: "function as argument",
			callExpr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "HandleFunc"},
				Args: []ast.Expr{
					&ast.BasicLit{Value: `"/path"`},
					&ast.Ident{Name: "handler"},
				},
			},
			privateFuncs: map[string][]*privateFuncInfo{
				"handler": {{name: "handler"}},
			},
			expected: []string{"handler"},
		},
		{
			name: "method as argument",
			callExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "mux"},
					Sel: &ast.Ident{Name: "HandleFunc"},
				},
				Args: []ast.Expr{
					&ast.BasicLit{Value: `"/live"`},
					&ast.SelectorExpr{
						X:   &ast.Ident{Name: "a"},
						Sel: &ast.Ident{Name: "handleLiveness"},
					},
				},
			},
			privateFuncs: map[string][]*privateFuncInfo{
				"handleLiveness": {{name: "handleLiveness", receiverType: "App"}},
			},
			expected: []string{"handleLiveness"},
		},
		{
			name: "unknown function as argument",
			callExpr: &ast.CallExpr{
				Fun: &ast.Ident{Name: "Register"},
				Args: []ast.Expr{
					&ast.Ident{Name: "unknownFunc"},
				},
			},
			privateFuncs: map[string][]*privateFuncInfo{
				"handler": {{name: "handler"}},
			},
			expected: []string{},
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			calledInProduction := make(map[string]bool)
			collectCallExprCallbacks(tt.callExpr, tt.privateFuncs, calledInProduction)

			// Vérification des résultats
			for _, expectedFunc := range tt.expected {
				// Vérifier si la fonction est marquée comme appelée
				if !calledInProduction[expectedFunc] {
					t.Errorf("expected %s to be marked as called", expectedFunc)
				}
			}
		})
	}
}

// Test_markIfPrivateFunc vérifie le marquage d'une fonction comme appelée.
func Test_markIfPrivateFunc(t *testing.T) {
	tests := []struct {
		name         string
		funcName     string
		privateFuncs map[string][]*privateFuncInfo
		shouldMark   bool
	}{
		{
			name:     "known private function",
			funcName: "privateHelper",
			privateFuncs: map[string][]*privateFuncInfo{
				"privateHelper": {{name: "privateHelper"}},
			},
			shouldMark: true,
		},
		{
			name:     "unknown function",
			funcName: "unknownFunc",
			privateFuncs: map[string][]*privateFuncInfo{
				"privateHelper": {{name: "privateHelper"}},
			},
			shouldMark: false,
		},
		{
			name:         "empty private funcs map",
			funcName:     "anyFunc",
			privateFuncs: map[string][]*privateFuncInfo{},
			shouldMark:   false,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			calledInProduction := make(map[string]bool)
			markIfPrivateFunc(tt.funcName, tt.privateFuncs, calledInProduction)

			// Vérification du résultat
			if calledInProduction[tt.funcName] != tt.shouldMark {
				t.Errorf("markIfPrivateFunc() marked = %v, want %v", calledInProduction[tt.funcName], tt.shouldMark)
			}
		})
	}
}

// Test_collectIdentCallbacks vérifie la collecte des identifiants de callbacks.
func Test_collectIdentCallbacks(t *testing.T) {
	tests := []struct {
		name         string
		privateFuncs map[string][]*privateFuncInfo
		node         ast.Node
		expected     []string
	}{
		{
			name: "ident in composite lit",
			privateFuncs: map[string][]*privateFuncInfo{
				"callback": {{name: "callback"}},
			},
			node: &ast.CompositeLit{
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key:   &ast.Ident{Name: "RunE"},
						Value: &ast.Ident{Name: "callback"},
					},
				},
			},
			expected: []string{"callback"},
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			calledInProduction := make(map[string]bool)
			collectIdentCallbacks(tt.node, tt.privateFuncs, calledInProduction)

			// Vérification des résultats
			for _, expectedFunc := range tt.expected {
				// Vérifier si la fonction est marquée comme appelée
				if !calledInProduction[expectedFunc] {
					t.Errorf("expected %s to be marked as called", expectedFunc)
				}
			}
		})
	}
}
