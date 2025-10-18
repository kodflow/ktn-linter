package ktn

import (
	"golang.org/x/tools/go/analysis"

	ktn_alloc "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/alloc"
	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
	ktn_data_structures "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/data_structures"
	ktn_error "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/error"
	ktn_func "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/func"
	ktn_goroutine "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/goroutine"
	ktn_interface "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/interface"
	ktn_method "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/method"
	ktn_mock "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/mock"
	ktn_ops "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/ops"
	ktn_package "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/package"
	ktn_pool "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/pool"
	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
	ktn_test "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/test"
	ktn_var "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/var"
)

// AllRules contient TOUTES les règles KTN organisées par catégorie.
//
// Total: ~64 règles réparties en 16 catégories
var AllRules = struct {
	Func           []*analysis.Analyzer // 9 règles
	Var            []*analysis.Analyzer // 9 règles
	Struct         []*analysis.Analyzer // 4 règles
	Interface      []*analysis.Analyzer // 6 règles
	Const          []*analysis.Analyzer // 4 règles
	Error          []*analysis.Analyzer // 1 règle
	Test           []*analysis.Analyzer // 4 règles
	Alloc          []*analysis.Analyzer // 3 règles
	Goroutine      []*analysis.Analyzer // 2 règles
	Pool           []*analysis.Analyzer // 1 règle
	Mock           []*analysis.Analyzer // 2 règles
	Method         []*analysis.Analyzer // 1 règle
	Package        []*analysis.Analyzer // 1 règle
	ControlFlow    []*analysis.Analyzer // 7 règles
	DataStructures []*analysis.Analyzer // 3 règles
	Ops            []*analysis.Analyzer // 8 règles
}{
	Func:           ktn_func.GetRules(),
	Var:            ktn_var.GetRules(),
	Struct:         ktn_struct.GetRules(),
	Interface:      ktn_interface.GetRules(),
	Const:          ktn_const.GetRules(),
	Error:          ktn_error.GetRules(),
	Test:           ktn_test.GetRules(),
	Alloc:          ktn_alloc.GetRules(),
	Goroutine:      ktn_goroutine.GetRules(),
	Pool:           ktn_pool.GetRules(),
	Mock:           ktn_mock.GetRules(),
	Method:         ktn_method.GetRules(),
	Package:        ktn_package.GetRules(),
	ControlFlow:    ktn_control_flow.GetRules(),
	DataStructures: ktn_data_structures.GetRules(),
	Ops:            ktn_ops.GetRules(),
}

// GetAllRules retourne toutes les règles KTN en une seule liste plate.
//
// Returns:
//   - []*analysis.Analyzer: liste de toutes les règles KTN (~64 règles)
func GetAllRules() []*analysis.Analyzer {
	var all []*analysis.Analyzer

	all = append(all, AllRules.Func...)
	all = append(all, AllRules.Var...)
	all = append(all, AllRules.Struct...)
	all = append(all, AllRules.Interface...)
	all = append(all, AllRules.Const...)
	all = append(all, AllRules.Error...)
	all = append(all, AllRules.Test...)
	all = append(all, AllRules.Alloc...)
	all = append(all, AllRules.Goroutine...)
	all = append(all, AllRules.Pool...)
	all = append(all, AllRules.Mock...)
	all = append(all, AllRules.Method...)
	all = append(all, AllRules.Package...)
	all = append(all, AllRules.ControlFlow...)
	all = append(all, AllRules.DataStructures...)
	all = append(all, AllRules.Ops...)

	// Early return from function.
	return all
}

// GetRulesByCategory retourne les règles d'une catégorie spécifique.
//
// Params:
//   - category: nom de la catégorie (func, var, struct, etc.)
//
// Returns:
//   - []*analysis.Analyzer: règles de la catégorie demandée
func GetRulesByCategory(category string) []*analysis.Analyzer {
	switch category {
	case "func":
		// Early return from function.
		return AllRules.Func
	case "var":
		// Early return from function.
		return AllRules.Var
	case "struct":
		// Early return from function.
		return AllRules.Struct
	case "interface":
		// Early return from function.
		return AllRules.Interface
	case "const":
		// Early return from function.
		return AllRules.Const
	case "error":
		// Early return from function.
		return AllRules.Error
	case "test":
		// Early return from function.
		return AllRules.Test
	case "alloc":
		// Early return from function.
		return AllRules.Alloc
	case "goroutine":
		// Early return from function.
		return AllRules.Goroutine
	case "pool":
		// Early return from function.
		return AllRules.Pool
	case "mock":
		// Early return from function.
		return AllRules.Mock
	case "method":
		// Early return from function.
		return AllRules.Method
	case "package":
		// Early return from function.
		return AllRules.Package
	case "control_flow":
		// Early return from function.
		return AllRules.ControlFlow
	case "data_structures":
		// Early return from function.
		return AllRules.DataStructures
	case "ops":
		// Early return from function.
		return AllRules.Ops
	default:
		// Early return from function.
		return nil
	}
}
