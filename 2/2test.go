package main

import (
	"fmt"
	"testing"
)

func UnpackTest(t *testing.T) {
	testCase := []struct {
		input    string
		expected string
	}{
		{"a4bc2d5e", "aaaabccddddde"},
		{"abcd", "abcd"},
		{"45", ""},
		{"", ""},
	}

	for _, tc := range testCase {
		t.Run(tc.input, func(t *testing.T) {
			res := unpack(tc.input)
			if res != tc.expected {
				fmt.Printf("Wrong result %s for %s. Expected: %s", res, tc.input, tc.expected)
			}
		})
	}

}
