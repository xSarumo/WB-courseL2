package main

import (
	"reflect"
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{
			name:    "single field",
			input:   "1",
			want:    []int{1},
			wantErr: false,
		},
		{
			name:    "multiple fields",
			input:   "1,3,5",
			want:    []int{1, 3, 5},
			wantErr: false,
		},
		{
			name:    "field range",
			input:   "1-3",
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "combined fields and ranges",
			input:   "1,3-5,7",
			want:    []int{1, 3, 4, 5, 7},
			wantErr: false,
		},
		{
			name:    "invalid number",
			input:   "a",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid range format",
			input:   "1-3-5",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid range order",
			input:   "5-3",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		fields        []int
		delimiter     string
		separatedOnly bool
		want          string
	}{
		{
			name:          "basic tab-separated",
			line:          "a\tb\tc\td",
			fields:        []int{1, 3},
			delimiter:     "\t",
			separatedOnly: false,
			want:          "a\tc",
		},
		{
			name:          "custom delimiter",
			line:          "a,b,c,d",
			fields:        []int{2, 4},
			delimiter:     ",",
			separatedOnly: false,
			want:          "b,d",
		},
		{
			name:          "out of range fields",
			line:          "a\tb\tc",
			fields:        []int{1, 4},
			delimiter:     "\t",
			separatedOnly: false,
			want:          "a",
		},
		{
			name:          "separated only - with delimiter",
			line:          "a\tb\tc",
			fields:        []int{1, 2},
			delimiter:     "\t",
			separatedOnly: true,
			want:          "a\tb",
		},
		{
			name:          "separated only - without delimiter",
			line:          "abc",
			fields:        []int{1},
			delimiter:     "\t",
			separatedOnly: true,
			want:          "",
		},
		{
			name:          "all fields out of range",
			line:          "a,b",
			fields:        []int{3, 4},
			delimiter:     ",",
			separatedOnly: false,
			want:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := processLine(tt.line, tt.fields, tt.delimiter, tt.separatedOnly)
			if !ok {
				got = ""
			}
			if got != tt.want {
				t.Errorf("processLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
