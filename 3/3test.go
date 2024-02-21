package main

import (
	"fmt"
	"reflect"
	"testing"
)

func ConvertToDoubleArrayTest(t *testing.T) {
	testCase := []struct {
		input    []string
		expected [][]string
	}{
		{[]string{"hello", "world!"}, [][]string{{"h", "e", "l", "l", "o"}, {"w", "o", "r", "l", "d", "!"}}},
		{[]string{"привет", "мир!"}, [][]string{{"п", "р", "и", "в", "е", "т"}, {"м", "и", "р", "!"}}},
	}

	for _, tCase := range testCase {
		t.Run(tCase.input[0], func(t *testing.T) {
			res := convertToDoubleArray(tCase.input)
			if !reflect.DeepEqual(res, tCase.expected) {
				fmt.Printf("Wrong result %v for %v. Expected: %v", res, tCase.input, tCase.expected)
			}
		})
	}
}
