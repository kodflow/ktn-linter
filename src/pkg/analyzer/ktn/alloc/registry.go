package ktn_alloc
import "golang.org/x/tools/go/analysis"
var AllRules = []*analysis.Analyzer{Rule001, Rule002, Rule003}
func GetRules() []*analysis.Analyzer {return AllRules}
