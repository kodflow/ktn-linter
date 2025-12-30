// External tests for severity.go (black-box testing).
package severity_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/severity"
)

// TestLevel_String tests the String method of Level type.
func TestLevel_String(t *testing.T) {
	tests := []struct {
		name  string
		level severity.Level
		want  string
	}{
		{name: "INFO level", level: severity.SeverityInfo, want: "INFO"},
		{name: "WARNING level", level: severity.SeverityWarning, want: "WARNING"},
		{name: "ERROR level", level: severity.SeverityError, want: "ERROR"},
		{name: "Unknown level", level: severity.Level(999), want: "UNKNOWN"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetSeverity tests the GetSeverity function.
func TestGetSeverity(t *testing.T) {
	tests := []struct {
		name     string
		ruleCode string
		want     severity.Level
	}{
		{name: "Known rule KTN-VAR-001 is ERROR", ruleCode: "KTN-VAR-001", want: severity.SeverityError},
		{name: "Known rule KTN-VAR-003 is WARNING", ruleCode: "KTN-VAR-003", want: severity.SeverityWarning},
		{name: "Known rule KTN-CONST-002 is INFO", ruleCode: "KTN-CONST-002", want: severity.SeverityInfo},
		{name: "Unknown rule defaults to WARNING", ruleCode: "KTN-UNKNOWN-999", want: severity.SeverityWarning},
		{name: "Empty rule code defaults to WARNING", ruleCode: "", want: severity.SeverityWarning},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if got := severity.GetSeverity(tt.ruleCode); got != tt.want {
				t.Errorf("GetSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLevel_ColorCode tests the ColorCode method of Level type.
func TestLevel_ColorCode(t *testing.T) {
	tests := []struct {
		name  string
		level severity.Level
		want  string
	}{
		{name: "INFO color is blue", level: severity.SeverityInfo, want: "\033[34m"},
		{name: "WARNING color is yellow", level: severity.SeverityWarning, want: "\033[33m"},
		{name: "ERROR color is red", level: severity.SeverityError, want: "\033[31m"},
		{name: "Unknown level defaults to white", level: severity.Level(999), want: "\033[37m"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.ColorCode(); got != tt.want {
				t.Errorf("Level.ColorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLevel_Symbol tests the Symbol method of Level type.
func TestLevel_Symbol(t *testing.T) {
	tests := []struct {
		name  string
		level severity.Level
		want  string
	}{
		{name: "INFO symbol", level: severity.SeverityInfo, want: "ℹ"},
		{name: "WARNING symbol", level: severity.SeverityWarning, want: "⚠"},
		{name: "ERROR symbol", level: severity.SeverityError, want: "✖"},
		{name: "Unknown level symbol", level: severity.Level(999), want: "●"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Symbol(); got != tt.want {
				t.Errorf("Level.Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
