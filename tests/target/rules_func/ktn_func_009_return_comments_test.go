package rules_func_test

import (
	"github.com/kodflow/ktn-linter/tests/target/rules_func"

	"testing"
)

func TestFindMaxValueWithReturnComments(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		wantMax   int
		wantFound bool
	}{
		{
			name:      "empty slice",
			values:    []int{},
			wantMax:   0,
			wantFound: false,
		},
		{
			name:      "single value",
			values:    []int{42},
			wantMax:   42,
			wantFound: true,
		},
		{
			name:      "multiple values",
			values:    []int{3, 7, 2, 9, 1},
			wantMax:   9,
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMax, gotFound := rules_func.FindMaxValueWithReturnComments(tt.values)
			if gotMax != tt.wantMax {
				t.Errorf("rules_func.FindMaxValueWithReturnComments() max = %v, want %v", gotMax, tt.wantMax)
			}
			if gotFound != tt.wantFound {
				t.Errorf("rules_func.FindMaxValueWithReturnComments() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestValidateInputWithReturnComments(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{
			name:    "valid value",
			value:   50,
			wantErr: false,
		},
		{
			name:    "negative value",
			value:   -1,
			wantErr: true,
		},
		{
			name:    "value too large",
			value:   101,
			wantErr: true,
		},
		{
			name:    "boundary min",
			value:   0,
			wantErr: false,
		},
		{
			name:    "boundary max",
			value:   100,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules_func.ValidateInputWithReturnComments(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("rules_func.ValidateInputWithReturnComments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDivideNumbersWithReturnComments(t *testing.T) {
	tests := []struct {
		name    string
		a       int
		b       int
		want    float64
		wantErr bool
	}{
		{
			name:    "normal division",
			a:       10,
			b:       2,
			want:    5.0,
			wantErr: false,
		},
		{
			name:    "division by zero",
			a:       10,
			b:       0,
			want:    0,
			wantErr: true,
		},
		{
			name:    "result with decimal",
			a:       7,
			b:       2,
			want:    3.5,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rules_func.DivideNumbersWithReturnComments(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("rules_func.DivideNumbersWithReturnComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("rules_func.DivideNumbersWithReturnComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessWithMultipleExitsWithComments(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		want    string
		wantErr bool
	}{
		{
			name:    "negative value",
			value:   -5,
			want:    "",
			wantErr: true,
		},
		{
			name:    "zero",
			value:   0,
			want:    "zero",
			wantErr: false,
		},
		{
			name:    "small value",
			value:   5,
			want:    "small",
			wantErr: false,
		},
		{
			name:    "medium value",
			value:   50,
			want:    "medium",
			wantErr: false,
		},
		{
			name:    "large value",
			value:   150,
			want:    "large",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rules_func.ProcessWithMultipleExitsWithComments(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("rules_func.ProcessWithMultipleExitsWithComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("rules_func.ProcessWithMultipleExitsWithComments() = %v, want %v", got, tt.want)
			}
		})
	}
}
