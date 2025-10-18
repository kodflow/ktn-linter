package ktn_goroutine
import "golang.org/x/tools/go/analysis"
var AllRules = []*analysis.Analyzer{Rule001, Rule002}
func GetRules() []*analysis.Analyzer {return AllRules}
