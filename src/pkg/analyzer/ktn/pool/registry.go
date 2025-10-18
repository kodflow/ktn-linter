package ktn_pool
import "golang.org/x/tools/go/analysis"
var AllRules = []*analysis.Analyzer{Rule001}
func GetRules() []*analysis.Analyzer {return AllRules}
