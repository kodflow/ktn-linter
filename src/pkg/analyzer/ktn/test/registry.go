package ktn_test

import "golang.org/x/tools/go/analysis"

// AllRules contains all analyzer rules for this category.
var AllRules []*analysis.Analyzer = []*analysis.Analyzer{
	Rule001, Rule002, Rule003, Rule004,
}

// GetRules returns all analyzer rules for this category.
func GetRules() []*analysis.Analyzer {
	// Early return from function.
	return AllRules
}
