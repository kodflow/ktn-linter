package shared

import (
	"testing"
)

// TestCommentConstants teste les constantes utilisées dans la vérification des commentaires.
//
// Params:
//   - t: instance de testing pour rapporter les erreurs
func TestCommentConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{
			name:     "commentPrefixLength should be 2",
			value:    commentPrefixLength,
			expected: 2,
		},
		{
			name:     "wantMinLength should be 8",
			value:    wantMinLength,
			expected: 8,
		},
		{
			name:     "wantKeywordLength should be 4",
			value:    wantKeywordLength,
			expected: 4,
		},
		{
			name:     "wantWithSpaceLength should be 5",
			value:    wantWithSpaceLength,
			expected: 5,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérification de la valeur de la constante
			if tt.value != tt.expected {
				t.Errorf("constant = %d, want %d", tt.value, tt.expected)
			}
		})
	}
}

// TestCommentParsing teste le parsing des commentaires avec les constantes définies.
//
// Params:
//   - t: instance de testing pour rapporter les erreurs
func TestCommentParsing(t *testing.T) {
	tests := []struct {
		name          string
		commentText   string
		shouldBeWant  bool
		minLengthMet  bool
	}{
		{
			name:         "short comment",
			commentText:  "// ok",
			shouldBeWant: false,
			minLengthMet: false,
		},
		{
			name:         "exact min length want directive",
			commentText:  "// want ",
			shouldBeWant: true,
			minLengthMet: true,
		},
		{
			name:         "block comment want",
			commentText:  "/* want */",
			shouldBeWant: true,
			minLengthMet: true,
		},
	}

	// Itération sur les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Vérification de la longueur minimale
			if (len(tt.commentText) >= wantMinLength) != tt.minLengthMet {
				t.Errorf("minLengthMet mismatch for %q: got %v, want %v",
					tt.commentText, len(tt.commentText) >= wantMinLength, tt.minLengthMet)
			}

			// Test extraction du contenu après le préfixe
			if len(tt.commentText) >= commentPrefixLength {
				content := tt.commentText[commentPrefixLength:]
				// Vérification que le contenu est extrait correctement
				if content == "" && len(tt.commentText) > commentPrefixLength {
					t.Errorf("content extraction failed for %q", tt.commentText)
				}
			}
		})
	}
}
