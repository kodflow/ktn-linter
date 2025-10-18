package ktn_data_structures

import "golang.org/x/tools/go/analysis"

// AllRules contains all analyzer rules for this category.
var AllRules []*analysis.Analyzer = []*analysis.Analyzer{RuleArray001, RuleMap001, RuleSlice001}

// GetRules returns all analyzer rules for this category.
func GetRules() []*analysis.Analyzer { return AllRules }
