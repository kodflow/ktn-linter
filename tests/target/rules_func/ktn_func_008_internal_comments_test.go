package rules_func_test

import "testing"

func TestComplexCalculationWithInternalComments(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		want    int
		wantErr bool
	}{
		{
			name:    "negative value",
			value:   -1,
			want:    0,
			wantErr: true,
		},
		{
			name:    "zero value",
			value:   0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "positive value",
			value:   10,
			want:    -10,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := complexCalculationWithInternalComments(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("complexCalculationWithInternalComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("complexCalculationWithInternalComments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessDataWithComments(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want []int
	}{
		{
			name: "empty data",
			data: []int{},
			want: []int{},
		},
		{
			name: "mixed values",
			data: []int{10, 5, 4, 3, -1, 0},
			want: []int{20, 15, 2, 3},
		},
		{
			name: "all negative",
			data: []int{-1, -2, -3},
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processDataWithComments(tt.data)
			if len(got) != len(tt.want) {
				t.Errorf("processDataWithComments() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("processDataWithComments()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
