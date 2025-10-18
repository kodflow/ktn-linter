package ktn_mock

import "golang.org/x/tools/go/analysis"

// AllRules contains all analyzer rules for this category.
var AllRules []*analysis.Analyzer = []*analysis.Analyzer{Rule001, Rule002}

// GetRules returns all analyzer rules for this category.
func GetRules() []*analysis.Analyzer { return AllRules }
