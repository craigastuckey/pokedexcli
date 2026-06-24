package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"   ", []string{}},
		{"hello world", []string{"hello", "world"}},
		{"  hello   world  ", []string{"hello", "world"}},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		for i, word := range actual {
			if word != c.expected[i] {
				t.Errorf("cleanInput(%q) = %v, want %v", c.input, actual, c.expected)
			}
		}
	}
}
