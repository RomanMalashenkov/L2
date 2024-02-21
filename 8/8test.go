package main

import (
	"fmt"
	"testing"
)

func RouteUpTest(t *testing.T) {
	testCase := []struct {
		input, expected string
	}{
		{"C:\\Windows\\system32", "C:\\Windows"},
		{"C:\\Users\\rykov\\GolandProjects\\L2\\shell", "C:\\Users\\rykov\\GolandProjects\\L2"},
		{"C:\\Users", "C:\\"},
	}

	for _, tCase := range testCase {
		t.Run(tCase.input, func(t *testing.T) {
			res := routeUp(tCase.input)
			if res != tCase.expected {
				fmt.Printf("Wrong result %s for %s. Expected: %s", res, tCase.input, tCase.expected)
			}
		})
	}
}
