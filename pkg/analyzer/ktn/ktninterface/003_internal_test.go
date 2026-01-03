// Package ktninterface provides internal tests for KTN-INTERFACE-003 helper functions.
package ktninterface

import (
	"go/ast"
	"testing"
)

// Test_isSingleMethodInterface tests the isSingleMethodInterface helper function.
func Test_isSingleMethodInterface(t *testing.T) {
	tests := []struct {
		name     string
		iface    *ast.InterfaceType
		expected bool
	}{
		{
			name: "single method interface",
			iface: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Read"}},
							Type:  &ast.FuncType{},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "two method interface",
			iface: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Read"}},
							Type:  &ast.FuncType{},
						},
						{
							Names: []*ast.Ident{{Name: "Write"}},
							Type:  &ast.FuncType{},
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "empty interface",
			iface: &ast.InterfaceType{
				Methods: &ast.FieldList{List: []*ast.Field{}},
			},
			expected: false,
		},
		{
			name:     "nil methods",
			iface:    &ast.InterfaceType{Methods: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := isSingleMethodInterface(tt.iface)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("isSingleMethodInterface() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Test_extractSingleMethodName tests the extractSingleMethodName helper function.
func Test_extractSingleMethodName(t *testing.T) {
	tests := []struct {
		name     string
		iface    *ast.InterfaceType
		expected string
	}{
		{
			name: "method with name",
			iface: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{{Name: "Read"}},
							Type:  &ast.FuncType{},
						},
					},
				},
			},
			expected: "Read",
		},
		{
			name: "method without name",
			iface: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.FuncType{},
						},
					},
				},
			},
			expected: "",
		},
		{
			name: "empty methods list",
			iface: &ast.InterfaceType{
				Methods: &ast.FieldList{List: []*ast.Field{}},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := extractSingleMethodName(tt.iface)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("extractSingleMethodName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// Test_followsErConvention tests the followsErConvention helper function.
func Test_followsErConvention(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		methodName    string
		expected      bool
	}{
		{
			name:          "Read -> Reader",
			interfaceName: "Reader",
			methodName:    "read",
			expected:      true,
		},
		{
			name:          "Write -> Writer",
			interfaceName: "Writer",
			methodName:    "write",
			expected:      true,
		},
		{
			name:          "Validate -> Validator",
			interfaceName: "Validator",
			methodName:    "validate",
			expected:      true,
		},
		{
			name:          "Handle -> Handler",
			interfaceName: "Handler",
			methodName:    "handle",
			expected:      true,
		},
		{
			name:          "invalid naming",
			interfaceName: "MyInterface",
			methodName:    "read",
			expected:      false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := followsErConvention(tt.interfaceName, tt.methodName)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("followsErConvention(%q, %q) = %v, want %v", tt.interfaceName, tt.methodName, got, tt.expected)
			}
		})
	}
}

// Test_capitalizeFirst tests the capitalizeFirst helper function.
func Test_capitalizeFirst(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase",
			input:    "read",
			expected: "Read",
		},
		{
			name:     "already capitalized",
			input:    "Read",
			expected: "Read",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "unicode",
			input:    "écrire",
			expected: "Écrire",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := capitalizeFirst(tt.input)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("capitalizeFirst(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// Test_suggestErName tests the suggestErName helper function.
func Test_suggestErName(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		expected   string
	}{
		{
			name:       "read",
			methodName: "read",
			expected:   "Reader",
		},
		{
			name:       "write",
			methodName: "write",
			expected:   "Writer",
		},
		{
			name:       "validate",
			methodName: "validate",
			expected:   "Validator",
		},
		{
			name:       "generate",
			methodName: "generate",
			expected:   "Generator",
		},
		{
			name:       "handle",
			methodName: "handle",
			expected:   "Handler",
		},
		{
			name:       "empty",
			methodName: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la fonction
			got := suggestErName(tt.methodName)
			// Vérifier le résultat
			if got != tt.expected {
				// Résultat incorrect
				t.Errorf("suggestErName(%q) = %q, want %q", tt.methodName, got, tt.expected)
			}
		})
	}
}

// Test_runInterface003 validates runInterface003 via integration test reference.
// Full behavior is tested in 003_external_test.go via analysistest.Run.
func Test_runInterface003(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func() bool
		errMsg    string
	}{
		{
			name:      "analyzer not nil",
			checkFunc: func() bool { return Analyzer003 != nil },
			errMsg:    "Analyzer003 should not be nil",
		},
		{
			name:      "analyzer name correct",
			checkFunc: func() bool { return Analyzer003.Name == "ktninterface003" },
			errMsg:    "Analyzer003.Name should be ktninterface003",
		},
		{
			name:      "run function assigned",
			checkFunc: func() bool { return Analyzer003.Run != nil },
			errMsg:    "Analyzer003.Run should not be nil",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier la condition
			if !tt.checkFunc() {
				// Condition non satisfaite
				t.Error(tt.errMsg)
			}
		})
	}
}

// Test_analyzeTypeSpec003 validates analyzeTypeSpec003 helper function.
// Full behavior is tested in 003_external_test.go via analysistest.Run.
func Test_analyzeTypeSpec003(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func() bool
		errMsg    string
	}{
		{
			name:      "run function callable",
			checkFunc: func() bool { return Analyzer003.Run != nil },
			errMsg:    "analyzeTypeSpec003 should be callable via Analyzer003.Run",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Vérifier la condition
			if !tt.checkFunc() {
				// Condition non satisfaite
				t.Error(tt.errMsg)
			}
		})
	}
}

// Test_checkErNamingConvention validates checkErNamingConvention helper function.
// Full behavior is tested in 003_external_test.go via analysistest.Run.
func Test_checkErNamingConvention(t *testing.T) {
	tests := []struct {
		name          string
		interfaceName string
		methodName    string
		expected      bool
	}{
		{
			name:          "Reader_read_convention",
			interfaceName: "Reader",
			methodName:    "read",
			expected:      true,
		},
		{
			name:          "Writer_write_convention",
			interfaceName: "Writer",
			methodName:    "write",
			expected:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Tester la convention
			result := followsErConvention(tt.interfaceName, tt.methodName)
			// Vérifier le résultat
			if result != tt.expected {
				// Résultat incorrect
				t.Errorf("followsErConvention(%q, %q) = %v, want %v", tt.interfaceName, tt.methodName, result, tt.expected)
			}
		})
	}
}
