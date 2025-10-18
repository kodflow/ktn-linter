package ktn_test

import "golang.org/x/tools/go/analysis"

var AllRules = []*analysis.Analyzer{
	Rule001, Rule002, Rule003, Rule004,
}

func GetRules() []*analysis.Analyzer {
	return AllRules
}
