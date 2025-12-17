// Package rules provides rule information extraction and formatting utilities.
package rules

import (
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/pkg/analyzer/ktn"
)

const (
	// ktnPrefix is the standard prefix for rule codes.
	ktnPrefix string = "KTN-"
	// codePartsCount is the expected number of parts after KTN- prefix.
	codePartsCount int = 2
)

// RuleInfo contains complete information about a KTN rule.
// It includes the rule code, category, analyzer name, description and example.
type RuleInfo struct {
	Code        string // KTN-FUNC-001
	Category    string // func
	Name        string // ktnfunc001
	Description string // Short description
	GoodExample string // Content from good.go
}

// RulesOutput is the complete output structure for the rules command.
// It aggregates all rules with metadata for display purposes.
type RulesOutput struct {
	TotalCount int        // Total number of rules
	Categories []string   // Available categories
	Rules      []RuleInfo // All rules (filtered if requested)
}

// ExtractRuleCode extracts the KTN code from an analyzer Doc field.
//
// Params:
//   - doc: analyzer Doc field (format: "KTN-XXX-YYY: description")
//
// Returns:
//   - string: rule code (e.g., "KTN-FUNC-001") or empty if invalid format
func ExtractRuleCode(doc string) string {
	// Check prefix
	if !strings.HasPrefix(doc, ktnPrefix) {
		// Not a KTN rule
		return ""
	}

	// Use strings.Cut to split on colon
	code, _, found := strings.Cut(doc, ":")
	// Check if colon was found
	if !found {
		// No colon found - try to extract code anyway
		return extractCodeWithoutColon(doc)
	}

	// Validate format
	if !isValidRuleCode(code) {
		// Invalid code format
		return ""
	}

	// Return validated code
	return code
}

// extractCodeWithoutColon handles Doc fields without colon separator.
//
// Params:
//   - doc: analyzer Doc field without colon
//
// Returns:
//   - string: rule code or empty if invalid
func extractCodeWithoutColon(doc string) string {
	// Split by space to get first token
	parts := strings.Fields(doc)
	// Check if we have at least one part
	if len(parts) == 0 {
		// Empty doc
		return ""
	}

	// Check if first part looks like a rule code
	firstPart := parts[0]
	// Validate format
	if isValidRuleCode(firstPart) {
		// Return the code
		return firstPart
	}

	// Invalid format
	return ""
}

// isValidRuleCode checks if a string is a valid KTN rule code format.
//
// Params:
//   - code: potential rule code to validate
//
// Returns:
//   - bool: true if valid format (KTN-XXX-YYY)
func isValidRuleCode(code string) bool {
	// Must start with KTN-
	if !strings.HasPrefix(code, ktnPrefix) {
		// Invalid prefix
		return false
	}

	// Remove prefix and check parts
	rest := code[len(ktnPrefix):]
	parts := strings.Split(rest, "-")

	// Must have exactly 2 parts: CATEGORY and NUMBER
	return len(parts) == codePartsCount
}

// ExtractDescription extracts the description from an analyzer Doc field.
//
// Params:
//   - doc: analyzer Doc field (format: "KTN-XXX-YYY: description")
//
// Returns:
//   - string: description part after colon (trimmed)
func ExtractDescription(doc string) string {
	// Use strings.Cut to split on colon
	_, description, found := strings.Cut(doc, ":")
	// Check if colon was found and has content
	if !found || description == "" {
		// No description available - return full doc
		return doc
	}

	// Return trimmed description
	return strings.TrimSpace(description)
}

// ExtractCategory extracts category name from rule code.
//
// Params:
//   - code: rule code (e.g., "KTN-FUNC-001")
//
// Returns:
//   - string: category in lowercase (e.g., "func") or empty if invalid
func ExtractCategory(code string) string {
	// Check prefix
	if !strings.HasPrefix(code, ktnPrefix) {
		// Invalid code
		return ""
	}

	// Remove prefix
	rest := code[len(ktnPrefix):]
	// Split by dash
	parts := strings.Split(rest, "-")
	// Check parts count
	if len(parts) != codePartsCount {
		// Invalid format
		return ""
	}

	// Return category in lowercase
	return strings.ToLower(parts[0])
}

// GetAllRuleInfos extracts information from all available KTN rules.
//
// Returns:
//   - []RuleInfo: list of rule information (sorted by code)
func GetAllRuleInfos() []RuleInfo {
	// Get all analyzers
	analyzers := ktn.GetAllRules()
	// Convert to RuleInfo slice
	return analyzersToRuleInfos(analyzers)
}

// GetRuleInfosByCategory extracts information for rules in a category.
//
// Params:
//   - category: category name (e.g., "func", "var")
//
// Returns:
//   - []RuleInfo: list of rule information for the category
func GetRuleInfosByCategory(category string) []RuleInfo {
	// Get analyzers for category
	analyzers := ktn.GetRulesByCategory(category)
	// Convert to RuleInfo slice
	return analyzersToRuleInfos(analyzers)
}

// GetRuleInfoByCode extracts information for a single rule.
//
// Params:
//   - code: rule code (e.g., "KTN-FUNC-001")
//
// Returns:
//   - *RuleInfo: rule information or nil if not found
func GetRuleInfoByCode(code string) *RuleInfo {
	// Get analyzer by code
	analyzer := ktn.GetRuleByCode(code)
	// Check if found
	if analyzer == nil {
		// Rule not found
		return nil
	}

	// Convert to RuleInfo
	info := analyzerToRuleInfo(analyzer)

	// Return pointer
	return &info
}

// analyzersToRuleInfos converts analyzers slice to RuleInfo slice.
//
// Params:
//   - analyzers: slice of analyzers to convert
//
// Returns:
//   - []RuleInfo: converted and sorted rule information
func analyzersToRuleInfos(analyzers []*analysis.Analyzer) []RuleInfo {
	// Initialize result slice
	infos := make([]RuleInfo, 0, len(analyzers))

	// Convert each analyzer
	for _, a := range analyzers {
		// Extract info
		info := analyzerToRuleInfo(a)
		// Skip if not a valid KTN rule
		if info.Code == "" {
			// Not a KTN rule (e.g., modernize)
			continue
		}
		// Append to result
		infos = append(infos, info)
	}

	// Sort by code
	sort.Slice(infos, func(i, j int) bool {
		// Compare codes alphabetically
		return infos[i].Code < infos[j].Code
	})

	// Return sorted slice
	return infos
}

// analyzerToRuleInfo converts a single analyzer to RuleInfo.
//
// Params:
//   - a: analyzer to convert
//
// Returns:
//   - RuleInfo: extracted rule information
func analyzerToRuleInfo(a *analysis.Analyzer) RuleInfo {
	// Extract code from Doc
	code := ExtractRuleCode(a.Doc)
	// Extract description from Doc
	description := ExtractDescription(a.Doc)
	// Extract category from code
	category := ExtractCategory(code)

	// Build RuleInfo
	return RuleInfo{
		Code:        code,
		Category:    category,
		Name:        a.Name,
		Description: description,
		GoodExample: "", // Loaded separately if needed
	}
}

// GetCategories returns all available category names.
//
// Returns:
//   - []string: sorted list of category names
func GetCategories() []string {
	// Define available categories
	categories := []string{
		"api",
		"comment",
		"const",
		"func",
		"interface",
		"return",
		"struct",
		"test",
		"var",
	}

	// Already sorted alphabetically
	return categories
}
