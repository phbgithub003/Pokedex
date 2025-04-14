package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "pokemon go",
			expected: []string{"pokemon", "go"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "one\ttwo\nthree",
			expected: []string{"one", "two", "three"},
		},
		{
			input:    "pikachu   charizard bulbasaur",
			expected: []string{"pikachu", "charizard", "bulbasaur"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Check if lengths match
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) returned %d words, expected %d words",
				c.input, len(actual), len(c.expected))
			continue
		}

		// Check each word in the slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q)[%d] = %q, expected %q",
					c.input, i, word, expectedWord)
			}
		}

		// Alternative: use reflect.DeepEqual to compare slices
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("cleanInput(%q) = %v, expected %v",
				c.input, actual, c.expected)
		}
	}
}
