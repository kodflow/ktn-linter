package severity_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/severity"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name  string
		level severity.Level
		want  string
	}{
		{
			name:  "INFO level",
			level: severity.SEVERITY_INFO,
			want:  "INFO",
		},
		{
			name:  "WARNING level",
			level: severity.SEVERITY_WARNING,
			want:  "WARNING",
		},
		{
			name:  "ERROR level",
			level: severity.SEVERITY_ERROR,
			want:  "ERROR",
		},
		{
			name:  "Unknown level",
			level: severity.Level(999),
			want:  "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSeverity(t *testing.T) {
	tests := []struct {
		name     string
		ruleCode string
		want     severity.Level
	}{
		{
			name:     "Known rule KTN-VAR-001 is ERROR",
			ruleCode: "KTN-VAR-001",
			want:     severity.SEVERITY_ERROR,
		},
		{
			name:     "Known rule KTN-VAR-003 is WARNING",
			ruleCode: "KTN-VAR-003",
			want:     severity.SEVERITY_WARNING,
		},
		{
			name:     "Known rule KTN-CONST-002 is INFO (groupement)",
			ruleCode: "KTN-CONST-002",
			want:     severity.SEVERITY_INFO,
		},
		{
			name:     "Unknown rule defaults to WARNING",
			ruleCode: "KTN-UNKNOWN-999",
			want:     severity.SEVERITY_WARNING,
		},
		{
			name:     "Empty rule code defaults to WARNING",
			ruleCode: "",
			want:     severity.SEVERITY_WARNING,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := severity.GetSeverity(tt.ruleCode); got != tt.want {
				t.Errorf("GetSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_ColorCode(t *testing.T) {
	tests := []struct {
		name  string
		level severity.Level
		want  string
	}{
		{
			name:  "INFO color is blue",
			level: severity.SEVERITY_INFO,
			want:  "\033[34m",
		},
		{
			name:  "WARNING color is yellow",
			level: severity.SEVERITY_WARNING,
			want:  "\033[33m",
		},
		{
			name:  "ERROR color is red",
			level: severity.SEVERITY_ERROR,
			want:  "\033[31m",
		},
		{
			name:  "Unknown level defaults to white",
			level: severity.Level(999),
			want:  "\033[37m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.ColorCode(); got != tt.want {
				t.Errorf("Level.ColorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_Symbol(t *testing.T) {
	tests := []struct {
		name  string
		level severity.Level
		want  string
	}{
		{
			name:  "INFO symbol",
			level: severity.SEVERITY_INFO,
			want:  "ℹ",
		},
		{
			name:  "WARNING symbol",
			level: severity.SEVERITY_WARNING,
			want:  "⚠",
		},
		{
			name:  "ERROR symbol",
			level: severity.SEVERITY_ERROR,
			want:  "✖",
		},
		{
			name:  "Unknown level symbol",
			level: severity.Level(999),
			want:  "●",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Symbol(); got != tt.want {
				t.Errorf("Level.Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
