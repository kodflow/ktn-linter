package ktn_alloc

import "golang.org/x/tools/go/analysis"

// AllRules contains all analyzer rules for this category.
var AllRules []*analysis.Analyzer = []*analysis.Analyzer{Rule001, Rule002, Rule003}

// GetRules returns all analyzer rules for this category.
func GetRules() []*analysis.Analyzer { return AllRules }
