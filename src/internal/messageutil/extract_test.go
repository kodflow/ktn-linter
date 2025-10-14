package messageutil_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/src/internal/messageutil"
)

func TestExtractCode(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "KTN-CONST-001",
			message:  "[KTN-CONST-001] Constante 'MaxValue' déclarée individuellement.",
			expected: "KTN-CONST-001",
		},
		{
			name:     "KTN-CONST-002",
			message:  "[KTN-CONST-002] Groupe de constantes sans commentaire de groupe.",
			expected: "KTN-CONST-002",
		},
		{
			name:     "KTN-VAR-001",
			message:  "[KTN-VAR-001] Variable 'count' déclarée individuellement.",
			expected: "KTN-VAR-001",
		},
		{
			name:     "KTN-VAR-008",
			message:  "[KTN-VAR-008] Variable 'max_size' contient un underscore.",
			expected: "KTN-VAR-008",
		},
		{
			name:     "no code",
			message:  "Constante déclarée sans code d'erreur",
			expected: "UNKNOWN",
		},
		{
			name:     "incomplete code",
			message:  "[KTN-CONST without closing bracket",
			expected: "UNKNOWN",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := messageutil.ExtractCode(tt.message)
			if result != tt.expected {
				t.Errorf("messageutil.ExtractCode(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}

func TestExtractMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple message",
			message:  "[KTN-CONST-001] Constante déclarée individuellement.",
			expected: "Constante déclarée individuellement.",
		},
		{
			name:     "multiline message",
			message:  "[KTN-CONST-001] Constante déclarée individuellement.\nExemple:\n  const (...)",
			expected: "Constante déclarée individuellement.",
		},
		{
			name:     "message with extra spaces",
			message:  "[KTN-CONST-002]   Groupe sans commentaire  ",
			expected: "Groupe sans commentaire",
		},
		{
			name:     "no bracket",
			message:  "Message sans code d'erreur",
			expected: "Message sans code d'erreur",
		},
		{
			name:     "empty after bracket",
			message:  "[KTN-CONST-001]",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := messageutil.ExtractMessage(tt.message)
			if result != tt.expected {
				t.Errorf("messageutil.ExtractMessage(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}

func TestExtractSuggestion(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name: "with example",
			message: `[KTN-CONST-001] Constante déclarée individuellement.
Exemple:
  const (
      MaxValue int = 100
  )`,
			expected: `const (
      MaxValue int = 100
  )`,
		},
		{
			name: "with empty lines",
			message: `[KTN-CONST-002] Groupe sans commentaire.
Exemple:

  // Description
  const (
      X int = 1
  )

`,
			expected: `// Description
  const (
      X int = 1
  )`,
		},
		{
			name:     "no example",
			message:  "[KTN-CONST-001] Constante déclarée individuellement.",
			expected: "",
		},
		{
			name:     "empty after Exemple:",
			message:  "[KTN-CONST-001] Constante.\nExemple:\n\n\n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := messageutil.ExtractSuggestion(tt.message)
			if result != tt.expected {
				t.Errorf("messageutil.ExtractSuggestion() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractConstName(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "single quote name",
			message:  "Constante 'MaxValue' déclarée individuellement.",
			expected: "MaxValue",
		},
		{
			name:     "variable name",
			message:  "Variable 'count' sans type explicite.",
			expected: "count",
		},
		{
			name:     "multiple quotes",
			message:  "Variable 'maxSize' de type 'int' sans commentaire.",
			expected: "maxSize",
		},
		{
			name:     "no quotes",
			message:  "Constante déclarée sans guillemets",
			expected: "MyConst",
		},
		{
			name:     "incomplete quote",
			message:  "Constante 'MaxValue sans fermeture",
			expected: "MyConst",
		},
		{
			name:     "empty quotes",
			message:  "Constante '' vide",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := messageutil.ExtractConstName(tt.message)
			if result != tt.expected {
				t.Errorf("messageutil.ExtractConstName(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}

func TestExtractType(t *testing.T) {
	tests := []struct {
		name       string
		suggestion string
		expected   string
	}{
		// Types de base
		{name: "int type", suggestion: "MaxValue int = 100", expected: "int"},
		{name: "string type", suggestion: "Name string = \"hello\"", expected: "string"},
		{name: "bool type", suggestion: "IsActive bool = true", expected: "bool"},

		// Types signés
		{name: "int8", suggestion: "Age int8 = 25", expected: "int8"},
		{name: "int16", suggestion: "Port int16 = 8080", expected: "int16"},
		{name: "int32", suggestion: "Count int32 = 1000", expected: "int32"},
		{name: "int64", suggestion: "Size int64 = 9999", expected: "int64"},

		// Types non signés
		{name: "uint", suggestion: "ID uint = 1", expected: "uint"},
		{name: "uint8", suggestion: "Level uint8 = 5", expected: "uint8"},
		{name: "uint16", suggestion: "Port uint16 = 443", expected: "uint16"},
		{name: "uint32", suggestion: "Flags uint32 = 255", expected: "uint32"},
		{name: "uint64", suggestion: "Timestamp uint64 = 123456", expected: "uint64"},

		// Types flottants
		{name: "float32", suggestion: "Rate float32 = 3.14", expected: "float32"},
		{name: "float64", suggestion: "Pi float64 = 3.14159", expected: "float64"},

		// Autres types
		{name: "byte", suggestion: "Data byte = 0xFF", expected: "byte"},
		{name: "rune", suggestion: "Char rune = 'A'", expected: "rune"},
		{name: "complex64", suggestion: "Z complex64 = 1+2i", expected: "complex64"},
		{name: "complex128", suggestion: "W complex128 = 3+4i", expected: "complex128"},

		// Cas spéciaux
		{name: "with <type>", suggestion: "Value <type> = ...", expected: "int"},
		{name: "no type", suggestion: "MaxValue = 100", expected: "int"},
		{name: "empty", suggestion: "", expected: "int"},

		// Types au début
		{name: "type at start", suggestion: "int MaxValue = 100", expected: "int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := messageutil.ExtractType(tt.suggestion)
			if result != tt.expected {
				t.Errorf("messageutil.ExtractType(%q) = %q, want %q", tt.suggestion, result, tt.expected)
			}
		})
	}
}

func TestExtractCode_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "multiple brackets with KTN",
			message:  "[INFO] [KTN-CONST-001] Message",
			expected: "KTN-CONST-001", // Cherche spécifiquement [KTN-
		},
		{
			name:     "nested brackets",
			message:  "[KTN-CONST-[001]] Message",
			expected: "KTN-CONST-[001",
		},
		{
			name:     "bracket at end",
			message:  "Message [KTN-CONST-001]",
			expected: "KTN-CONST-001",
		},
		{
			name:     "non-KTN bracket first",
			message:  "[INFO] No KTN code",
			expected: "UNKNOWN", // Pas de [KTN- trouvé
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := messageutil.ExtractCode(tt.message)
			if result != tt.expected {
				t.Errorf("messageutil.ExtractCode(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}

func TestExtractMessage_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "multiple brackets",
			message:  "[CODE1] [CODE2] Message",
			expected: "[CODE2] Message",
		},
		{
			name:     "unicode message",
			message:  "[KTN-CONST-001] Message avec des caractères unicode: éàü",
			expected: "Message avec des caractères unicode: éàü",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := messageutil.ExtractMessage(tt.message)
			if result != tt.expected {
				t.Errorf("messageutil.ExtractMessage(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}
