package ktn_control_flow_test

import (
	"go/ast"
	"testing"

	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
)

// TestContainsNode teste directement la fonction containsNode pour améliorer la couverture
// TestContainsNode_NilBlock tests the functionality of the corresponding implementation.
func TestContainsNode_NilBlock(t *testing.T) {
	// Test avec block nil
	target := &ast.Ident{Name: "test"}
	result := ktn_control_flow.ContainsNodeExported(nil, target)
	if result {
		t.Errorf("containsNode(nil, target) devrait retourner false, obtenu true")
	}
}

// TestContainsNode_Found tests the functionality of the corresponding implementation.
func TestContainsNode_Found(t *testing.T) {
	// Test avec nœud trouvé
	target := &ast.Ident{Name: "test"}
	block := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: target,
			},
		},
	}
	result := ktn_control_flow.ContainsNodeExported(block, target)
	if !result {
		t.Errorf("containsNode devrait retourner true quand le nœud est trouvé")
	}
}

// TestContainsNode_NotFound tests the functionality of the corresponding implementation.
func TestContainsNode_NotFound(t *testing.T) {
	// Test avec nœud non trouvé
	target := &ast.Ident{Name: "test"}
	other := &ast.Ident{Name: "other"}
	block := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ExprStmt{
				X: other,
			},
		},
	}
	result := ktn_control_flow.ContainsNodeExported(block, target)
	if result {
		t.Errorf("containsNode devrait retourner false quand le nœud n'est pas trouvé")
	}
}
