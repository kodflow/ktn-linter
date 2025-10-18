package ktn_ops

import "golang.org/x/tools/go/analysis"

// AllRules contains all analyzer rules for this category.
var AllRules = []*analysis.Analyzer{RuleChan001, RuleComp001, RuleConv001, RuleAssert001, RuleOp001, RulePointer001, RulePredecl001, RuleReturn001}

// GetRules returns all analyzer rules for this category.
func GetRules() []*analysis.Analyzer { return AllRules }
