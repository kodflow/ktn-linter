package ktn_control_flow
import "golang.org/x/tools/go/analysis"
var AllRules = []*analysis.Analyzer{RuleDefer001, RuleFor001, RuleIf001, RuleGoto001, RuleSwitch001, RuleFall001, RuleRange001}
func GetRules() []*analysis.Analyzer {return AllRules}
