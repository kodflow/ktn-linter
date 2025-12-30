package shared_test

import (
	"go/ast"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
)

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "test file",
			filename: "example_test.go",
			want:     true,
		},
		{
			name:     "regular file",
			filename: "example.go",
			want:     false,
		},
		{
			name:     "test file with path",
			filename: "/path/to/file_test.go",
			want:     true,
		},
		{
			name:     "non-go file",
			filename: "example.txt",
			want:     false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := shared.IsTestFile(tt.filename)
			if got != tt.want {
				t.Errorf("IsTestFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTestFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		want     bool
	}{
		{
			name:     "nil function",
			funcDecl: nil,
			want:     false,
		},
		{
			name: "test function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("TestExample"),
			},
			want: true,
		},
		{
			name: "benchmark function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("BenchmarkExample"),
			},
			want: true,
		},
		{
			name: "example function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("ExampleFoo"),
			},
			want: true,
		},
		{
			name: "fuzz function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("FuzzBar"),
			},
			want: true,
		},
		{
			name: "regular function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("ProcessData"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := shared.IsTestFunction(tt.funcDecl)
			if got != tt.want {
				t.Errorf("IsTestFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnitTestFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		want     bool
	}{
		{
			name:     "nil function",
			funcDecl: nil,
			want:     false,
		},
		{
			name: "test function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("TestExample"),
			},
			want: true,
		},
		{
			name: "benchmark function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("BenchmarkExample"),
			},
			want: false,
		},
		{
			name: "example function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("ExampleFoo"),
			},
			want: false,
		},
		{
			name: "fuzz function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("FuzzBar"),
			},
			want: false,
		},
		{
			name: "regular function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("ProcessData"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := shared.IsUnitTestFunction(tt.funcDecl)
			if got != tt.want {
				t.Errorf("IsUnitTestFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExportedFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcDecl *ast.FuncDecl
		want     bool
	}{
		{
			name:     "nil function",
			funcDecl: nil,
			want:     false,
		},
		{
			name: "exported function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("PublicFunc"),
			},
			want: true,
		},
		{
			name: "unexported function",
			funcDecl: &ast.FuncDecl{
				Name: ast.NewIdent("privateFunc"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := shared.IsExportedFunction(tt.funcDecl)
			if got != tt.want {
				t.Errorf("IsExportedFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}
