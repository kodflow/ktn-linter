package formatter_test

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/src/pkg/formatter"
	"golang.org/x/tools/go/analysis"
)

// TestNewFormatter teste NewFormatter.
//
// Params:
//   - t: instance de test
func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name       string
		aiMode     bool
		noColor    bool
		simpleMode bool
	}{
		{"default mode", false, false, false},
		{"AI mode", true, false, false},
		{"no color mode", false, true, false},
		{"simple mode", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := formatter.NewFormatter(&buf, tt.aiMode, tt.noColor, tt.simpleMode)

			if f == nil {
				t.Fatal("NewFormatter returned nil")
			}
		})
	}
}

// TestFormatterFormatSuccess teste Formatter Format Success.
//
// Params:
//   - t: instance de test
func TestFormatterFormatSuccess(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, false, false)

	f.Format(nil, nil)

	output := buf.String()
	if !strings.Contains(output, "No issues found") {
		t.Errorf("Expected success message, got: %s", output)
	}
}

// TestFormatterFormatHumanMode teste Formatter Format HumanMode.
//
// Params:
//   - t: instance de test
func TestFormatterFormatHumanMode(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 100)

	diagnostics := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "[KTN-CONST-001] Constante 'MaxValue' déclarée individuellement",
		},
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false) // noColor=true

	f.Format(fset, diagnostics)

	output := buf.String()
	if !strings.Contains(output, "KTN-LINTER REPORT") {
		t.Errorf("Expected header, got: %s", output)
	}
	if !strings.Contains(output, "test.go") {
		t.Errorf("Expected filename, got: %s", output)
	}
	if !strings.Contains(output, "KTN-CONST-001") {
		t.Errorf("Expected error code, got: %s", output)
	}
}

// TestFormatterFormatAIMode teste Formatter Format AIMode.
//
// Params:
//   - t: instance de test
func TestFormatterFormatAIMode(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 100)

	diagnostics := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "[KTN-CONST-001] Constante 'MaxValue' déclarée individuellement",
		},
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, true, true, false) // aiMode=true

	f.Format(fset, diagnostics)

	output := buf.String()
	if !strings.Contains(output, "# KTN-Linter Report (AI Mode)") {
		t.Errorf("Expected AI mode header, got: %s", output)
	}
	if !strings.Contains(output, "## File: test.go") {
		t.Errorf("Expected file section, got: %s", output)
	}
	if !strings.Contains(output, "**Code**: KTN-CONST-001") {
		t.Errorf("Expected code in markdown, got: %s", output)
	}
}

// TestFormatterFormatSimpleMode teste Formatter Format SimpleMode.
//
// Params:
//   - t: instance de test
func TestFormatterFormatSimpleMode(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 100)

	diagnostics := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "[KTN-CONST-001] Constante 'MaxValue' déclarée individuellement",
		},
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, true) // simpleMode=true

	f.Format(fset, diagnostics)

	output := buf.String()
	// Simple mode: file:line:column: [CODE] message
	if !strings.Contains(output, "test.go:") {
		t.Errorf("Expected simple format with filename, got: %s", output)
	}
	if !strings.Contains(output, "[KTN-CONST-001]") {
		t.Errorf("Expected error code in brackets, got: %s", output)
	}
	if !strings.Contains(output, "Constante") {
		t.Errorf("Expected message, got: %s", output)
	}
}

// TestFormatterFormatMultipleDiagnostics teste Formatter Format MultipleDiagnostics.
//
// Params:
//   - t: instance de test
func TestFormatterFormatMultipleDiagnostics(t *testing.T) {
	fset := token.NewFileSet()
	file1 := fset.AddFile("file1.go", 1, 100)
	file2 := fset.AddFile("file2.go", 102, 200)

	diagnostics := []analysis.Diagnostic{
		{Pos: file1.Pos(10), Message: "[KTN-CONST-001] Error 1"},
		{Pos: file2.Pos(110), Message: "[KTN-VAR-003] Error 2"},
		{Pos: file1.Pos(20), Message: "[KTN-CONST-002] Error 3"},
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false) // noColor=true

	f.Format(fset, diagnostics)

	output := buf.String()

	// Should contain both files
	if !strings.Contains(output, "file1.go") {
		t.Errorf("Expected file1.go in output")
	}
	if !strings.Contains(output, "file2.go") {
		t.Errorf("Expected file2.go in output")
	}

	// Should show total count
	if !strings.Contains(output, "3 issue") {
		t.Errorf("Expected 3 issues in output, got: %s", output)
	}
}

// TestFormatterFormatColorsDisabled teste Formatter Format ColorsDisabled.
//
// Params:
//   - t: instance de test
func TestFormatterFormatColorsDisabled(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 100)

	diagnostics := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "[KTN-CONST-001] Test error",
		},
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false) // noColor=true

	f.Format(fset, diagnostics)

	output := buf.String()

	// Should not contain ANSI escape codes
	if strings.Contains(output, "\033[") {
		t.Errorf("Expected no ANSI codes with noColor=true, got: %s", output)
	}
}

// TestFormatterFormatWithColors teste Formatter Format WithColors.
//
// Params:
//   - t: instance de test
func TestFormatterFormatWithColors(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 100)

	diagnostics := []analysis.Diagnostic{
		{
			Pos:     file.Pos(10),
			Message: "[KTN-CONST-001] Test error",
		},
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, false, false) // with colors

	f.Format(fset, diagnostics)

	output := buf.String()

	// Should contain ANSI escape codes for colors
	if !strings.Contains(output, "\033[") {
		t.Errorf("Expected ANSI codes with noColor=false")
	}
}

// TestFormatterFormatSimpleModeSorting teste Formatter Format SimpleMode Sorting.
//
// Params:
//   - t: instance de test
func TestFormatterFormatSimpleModeSorting(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 100)

	// Add diagnostics in non-sorted order
	diagnostics := []analysis.Diagnostic{
		{Pos: file.Pos(30), Message: "[KTN-CONST-003] Error at line 30"},
		{Pos: file.Pos(10), Message: "[KTN-CONST-001] Error at line 10"},
		{Pos: file.Pos(20), Message: "[KTN-CONST-002] Error at line 20"},
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, true) // simpleMode=true

	f.Format(fset, diagnostics)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should be sorted by line number
	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines, got %d", len(lines))
	}

	// Check that line 10 comes before line 20 which comes before line 30
	line10Idx := -1
	line20Idx := -1
	line30Idx := -1

	for i, line := range lines {
		if strings.Contains(line, "Error at line 10") {
			line10Idx = i
		}
		if strings.Contains(line, "Error at line 20") {
			line20Idx = i
		}
		if strings.Contains(line, "Error at line 30") {
			line30Idx = i
		}
	}

	if line10Idx == -1 || line20Idx == -1 || line30Idx == -1 {
		t.Fatalf("Not all errors found in output: %s", output)
	}

	if !(line10Idx < line20Idx && line20Idx < line30Idx) {
		t.Errorf("Expected errors to be sorted by line number, got order: %d, %d, %d",
			line10Idx, line20Idx, line30Idx)
	}
}

// allErrorCodesDiagnostics contient tous les codes d'erreur KTN pour les tests de couverture complète.
//
// Params:
//   - file: le fichier token pour créer les positions
//
// Returns:
//   - []analysis.Diagnostic: diagnostics de test couvrant tous les codes d'erreur
func allErrorCodesDiagnostics(file *token.File) []analysis.Diagnostic {
	return []analysis.Diagnostic{
		{Pos: file.Pos(10), Message: "[KTN-CONST-001] Constante 'MaxValue' déclarée individuellement.\nExemple:\n  const (\n      MaxValue int = ...\n  )"},
		{Pos: file.Pos(20), Message: "[KTN-CONST-002] Groupe de constantes sans commentaire de groupe."},
		{Pos: file.Pos(30), Message: "[KTN-CONST-003] Constante 'Timeout' sans commentaire individuel.\nExemple:\n  // Timeout décrit son rôle\n  Timeout int = ..."},
		{Pos: file.Pos(40), Message: "[KTN-CONST-004] Constante 'BufferSize' sans type explicite.\nExemple:\n  BufferSize int = ..."},
		{Pos: file.Pos(50), Message: "[KTN-VAR-001] Variable 'MaxConnections' déclarée individuellement.\nExemple:\n  var (\n      MaxConnections int = ...\n  )"},
		{Pos: file.Pos(60), Message: "[KTN-VAR-002] Groupe de variables sans commentaire de groupe."},
		{Pos: file.Pos(70), Message: "[KTN-VAR-005] Variable 'Pi' avec valeur littérale semble être une constante immuable.\nExemple:\n  const Pi float64 = ..."},
		{Pos: file.Pos(80), Message: "Diagnostic without code format"},
	}
}

// TestFormatterHumanModeWithColors teste le mode humain avec couleurs.
//
// Params:
//   - t: instance de test
func TestFormatterHumanModeWithColors(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 200)
	diagnostics := allErrorCodesDiagnostics(file)

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, false, false)
	f.Format(fset, diagnostics)

	output := buf.String()
	// Should contain all error codes
	for _, code := range []string{"KTN-CONST-001", "KTN-CONST-002", "KTN-CONST-003", "KTN-CONST-004", "KTN-VAR-001", "KTN-VAR-002", "KTN-VAR-005"} {
		if !strings.Contains(output, code) {
			t.Errorf("Expected output to contain %s", code)
		}
	}
	// Should contain ANSI codes for colors
	if !strings.Contains(output, "\033[") {
		t.Error("Expected ANSI color codes in output")
	}
}

// TestFormatterAIModeWithSuggestions teste le mode AI avec suggestions.
//
// Params:
//   - t: instance de test
func TestFormatterAIModeWithSuggestions(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 200)
	diagnostics := allErrorCodesDiagnostics(file)

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, true, false, false)
	f.Format(fset, diagnostics)

	output := buf.String()
	// Should contain AI mode markers
	if !strings.Contains(output, "# KTN-Linter Report (AI Mode)") {
		t.Error("Expected AI mode header")
	}
	if !strings.Contains(output, "**Instructions for AI**") {
		t.Error("Expected AI instructions")
	}
	// Should contain suggestions in code blocks
	if !strings.Contains(output, "```go") {
		t.Error("Expected code block with suggestions")
	}
}

// TestFormatterSimpleMode teste le mode simple.
//
// Params:
//   - t: instance de test
func TestFormatterSimpleMode(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, 200)
	diagnostics := allErrorCodesDiagnostics(file)

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, true)
	f.Format(fset, diagnostics)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(diagnostics) {
		t.Errorf("Expected %d lines, got %d", len(diagnostics), len(lines))
	}
	// Each line should have format: file:line:column: [CODE] message
	for _, line := range lines {
		if !strings.Contains(line, "test.go:") {
			t.Errorf("Expected line to contain filename, got: %s", line)
		}
	}
}

// TestFormatterMessageWithoutCode teste les messages sans code.
//
// Params:
//   - t: instance de test
func TestFormatterMessageWithoutCode(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("edge.go", 1, 100)

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false)

	diagnostics := []analysis.Diagnostic{
		{Pos: file.Pos(10), Message: "Error without code format"},
	}

	f.Format(fset, diagnostics)
	output := buf.String()

	if !strings.Contains(output, "UNKNOWN") {
		t.Error("Expected UNKNOWN code for message without KTN- format")
	}
}

// TestFormatterMessageWithoutClosingBracket teste les messages mal formés.
//
// Params:
//   - t: instance de test
func TestFormatterMessageWithoutClosingBracket(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("edge.go", 1, 100)

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false)

	diagnostics := []analysis.Diagnostic{
		{Pos: file.Pos(10), Message: "[KTN-TEST Error without closing bracket"},
	}

	f.Format(fset, diagnostics)
	output := buf.String()

	if !strings.Contains(output, "UNKNOWN") {
		t.Error("Expected UNKNOWN code for malformed message")
	}
}

// TestFormatterMessageWithoutNewline teste les messages sans exemple.
//
// Params:
//   - t: instance de test
func TestFormatterMessageWithoutNewline(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("edge.go", 1, 100)

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false)

	diagnostics := []analysis.Diagnostic{
		{Pos: file.Pos(10), Message: "[KTN-TEST-001] Single line message without example"},
	}

	f.Format(fset, diagnostics)
	output := buf.String()

	if !strings.Contains(output, "Single line message") {
		t.Error("Expected message to be extracted")
	}
}

// TestFormatterSuccessWithColors teste le succès avec couleurs.
//
// Params:
//   - t: instance de test
func TestFormatterSuccessWithColors(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, false, false) // with colors

	f.Format(nil, nil)
	output := buf.String()

	if !strings.Contains(output, "No issues found") {
		t.Error("Expected success message")
	}
	if !strings.Contains(output, "\033[") {
		t.Error("Expected ANSI codes for colored success message")
	}
}

// TestFormatterSuccessWithoutColors teste le succès sans couleurs.
//
// Params:
//   - t: instance de test
func TestFormatterSuccessWithoutColors(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false) // noColor

	f.Format(nil, nil)
	output := buf.String()

	if !strings.Contains(output, "No issues found") {
		t.Error("Expected success message")
	}
	if strings.Contains(output, "\033[") {
		t.Error("Should not contain ANSI codes with noColor=true")
	}
}

// TestFormatterTypeExtraction tests type extraction from suggestions
//
// Params:
//   - t: instance de test
func TestFormatterTypeExtraction(t *testing.T) {
	fset := token.NewFileSet()
	file := fset.AddFile("types.go", 1, 300)

	// Test various Go types in suggestions
	types := []string{
		"bool", "string", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "byte", "rune",
		"complex64", "complex128",
	}

	diagnostics := []analysis.Diagnostic{}
	for i, typ := range types {
		diagnostics = append(diagnostics, analysis.Diagnostic{
			Pos: file.Pos(10 + i*5),
			Message: fmt.Sprintf("[KTN-CONST-001] Constante 'Value%d' déclarée individuellement.\nExemple:\n  const (\n      Value%d %s = ...\n  )", i, i, typ),
		})
	}

	var buf bytes.Buffer
	f := formatter.NewFormatter(&buf, false, true, false)
	f.Format(fset, diagnostics)

	output := buf.String()
	// Verify all types are present in output
	for _, typ := range types {
		if !strings.Contains(output, typ) {
			t.Errorf("Expected output to contain type %s", typ)
		}
	}
}
