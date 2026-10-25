package api

import (
	"reflect"
	"testing"
)

func TestSplitParts(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"ICE123|12:30|Berlin|Hannover|0|Info", []string{"ICE123", "12:30", "Berlin", "Hannover", "0", "Info"}},
		{"ICE123|12:30|Berlin|Hannover|0|Info|A|B|C", []string{"ICE123", "12:30", "Berlin", "Hannover", "0", "Info|A|B|C"}},
		{"ICE123|12:30|Berlin", []string{"ICE123", "12:30", "Berlin", "", "", ""}},
		{"", []string{"", "", "", "", "", ""}},
	}
	for _, tt := range tests {
		if got := splitParts(tt.in, 6); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("splitParts(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
