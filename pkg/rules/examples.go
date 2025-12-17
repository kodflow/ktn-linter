// Package rules provides rule information extraction and formatting utilities.
package rules

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// testdataBasePath is the base path for testdata directories.
	testdataBasePath string = "pkg/analyzer/ktn"
	// testdataSuffix is the subdirectory structure for testdata.
	testdataSuffix string = "testdata/src"
	// goodFileName is the standard name for good example files.
	goodFileName string = "good.go"
)

// GetTestdataPath returns the path to testdata directory for a rule.
//
// Params:
//   - code: rule code (e.g., "KTN-FUNC-001")
//
// Returns:
//   - string: path to testdata directory (e.g., "pkg/analyzer/ktn/ktnfunc/testdata/src/func001")
//   - error: error if code format is invalid
func GetTestdataPath(code string) (string, error) {
	// Extract category and number from code
	category, number, err := parseRuleCode(code)
	// Check for parsing errors
	if err != nil {
		// Return error
		return "", err
	}

	// Build testdata path
	// Format: pkg/analyzer/ktn/ktn<category>/testdata/src/<category><number>
	path := filepath.Join(
		testdataBasePath,
		"ktn"+category,
		testdataSuffix,
		category+number,
	)

	// Return path
	return path, nil
}

// parseRuleCode parses a rule code into category and number components.
//
// Params:
//   - code: rule code (e.g., "KTN-FUNC-001")
//
// Returns:
//   - string: category in lowercase (e.g., "func")
//   - string: number part (e.g., "001")
//   - error: error if format is invalid
func parseRuleCode(code string) (string, string, error) {
	// Check prefix
	if !strings.HasPrefix(code, ktnPrefix) {
		// Invalid prefix
		return "", "", &InvalidCodeError{Code: code, Reason: "missing KTN- prefix"}
	}

	// Remove prefix
	rest := code[len(ktnPrefix):]
	// Split by dash
	parts := strings.Split(rest, "-")
	// Check parts count
	if len(parts) != codePartsCount {
		// Invalid format
		return "", "", &InvalidCodeError{Code: code, Reason: "expected format KTN-CATEGORY-NNN"}
	}

	// Extract parts
	category := strings.ToLower(parts[0])
	number := parts[1]

	// Return parts
	return category, number, nil
}

// InvalidCodeError represents an error for invalid rule codes.
// It provides details about which code was invalid and why.
type InvalidCodeError struct {
	Code   string // The invalid code
	Reason string // Why it's invalid
}

// NewInvalidCodeError creates a new InvalidCodeError instance.
//
// Params:
//   - code: the invalid code
//   - reason: explanation of why the code is invalid
//
// Returns:
//   - *InvalidCodeError: new error instance
func NewInvalidCodeError(code, reason string) *InvalidCodeError {
	// Return new instance
	return &InvalidCodeError{Code: code, Reason: reason}
}

// Error implements the error interface.
//
// Returns:
//   - string: formatted error message
func (e *InvalidCodeError) Error() string {
	// Format error message
	return "invalid rule code '" + e.Code + "': " + e.Reason
}

// LoadGoodExample loads the good.go content for a rule.
//
// Params:
//   - code: rule code (e.g., "KTN-FUNC-001")
//
// Returns:
//   - string: content of good.go file or empty if not found
func LoadGoodExample(code string) string {
	// Get testdata path
	testdataPath, err := GetTestdataPath(code)
	// Check for errors
	if err != nil {
		// Invalid code - return empty
		return ""
	}

	// Try to find project root
	rootPath := findProjectRoot()
	// Build full path to good.go
	goodPath := filepath.Join(rootPath, testdataPath, goodFileName)

	// Read file content
	content, err := os.ReadFile(goodPath)
	// Check for errors
	if err != nil {
		// File not found or read error
		return ""
	}

	// Return content as string
	return string(content)
}

// findProjectRoot attempts to find the project root directory.
//
// Returns:
//   - string: path to project root or current directory if not found
func findProjectRoot() string {
	// Try using caller information first
	_, filename, _, ok := runtime.Caller(0)
	// Check if successful
	if ok {
		// Go up from pkg/rules/examples.go to project root
		dir := filepath.Dir(filename)
		// Go up twice: rules -> pkg -> project root
		projectRoot := filepath.Join(dir, "..", "..")
		// Verify by checking for go.mod
		if fileExists(filepath.Join(projectRoot, "go.mod")) {
			// Found project root
			return projectRoot
		}
	}

	// Fallback: try current working directory
	cwd, err := os.Getwd()
	// Check for errors
	if err != nil {
		// Return empty as last resort
		return ""
	}

	// Return current directory
	return cwd
}

// fileExists checks if a file exists at the given path.
//
// Params:
//   - path: path to check
//
// Returns:
//   - bool: true if file exists
func fileExists(path string) bool {
	// Check file info
	_, err := os.Stat(path)
	// Return true if no error
	return err == nil
}

// LoadGoodExamples loads good.go content for multiple rules.
//
// Params:
//   - infos: slice of RuleInfo to enrich with examples
//
// Returns:
//   - []RuleInfo: enriched rules with GoodExample populated
func LoadGoodExamples(infos []RuleInfo) []RuleInfo {
	// Create result slice with capacity only (no length)
	result := make([]RuleInfo, 0, len(infos))

	// Process each rule
	for _, info := range infos {
		// Copy info
		enriched := info
		// Load good example if code is valid
		if info.Code != "" {
			// Load example content
			enriched.GoodExample = LoadGoodExample(info.Code)
		}
		// Append to result
		result = append(result, enriched)
	}

	// Return enriched slice
	return result
}
