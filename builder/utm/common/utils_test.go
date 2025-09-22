package common

import "testing"

func TestMajorMinorDriverVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"4.7.1", "4_7"},
		{"4.7.9", "4_7"},
		{"4.7.11", "4_7"},
		{"5.0.0", "5_0"},
		{"10.11.22", "10_11"},
		{"4.7", "4.7"},
		{"4", "4"},
		{"", ""},
		{" 4.7.2", "4_7"}, // leading space
		{"4.7.2 ", "4_7"}, // trailing space
		{"something", "something"},
		{"4.7.1.1", "4.7.1.1"}, // too many segments
		{"4.7.x", "4.7.x"},     // non-numeric patch
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual := MajorMinorDriverVersion(test.input)

			if actual != test.expected {
				t.Errorf("MajorMinorDriverVersion(%q) = %q; expected %q", test.input, actual, test.expected)
			}
		})
	}
}
