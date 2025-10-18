package ktn_data_structures
import "golang.org/x/tools/go/analysis"
var AllRules = []*analysis.Analyzer{RuleArray001, RuleMap001, RuleSlice001}
func GetRules() []*analysis.Analyzer {return AllRules}
