package main

import "testing"

func TestIsYesAnswer(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"y", true},
		{"Y", true},
		{"yes", true},
		{"YES", true},
		{"Yes", true},
		{"n", false},
		{"no", false},
		{"", false},
		{"yeah", false},
	}

	for _, tt := range tests {
		if got := isYesAnswer(tt.input); got != tt.want {
			t.Errorf("isYesAnswer(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
