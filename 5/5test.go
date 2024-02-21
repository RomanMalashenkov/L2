package main

import (
	"fmt"
	"testing"
)

type testInput struct {
	word     string
	bothRegs bool
	revert   bool
}

func GenerateRegexTest(t *testing.T) {

	ti1 := testInput{
		"NATS",
		true,
		true,
	}

	ti2 := testInput{
		"NATS",
		false,
		true,
	}

	ti3 := testInput{
		"NATS",
		true,
		false,
	}

	ti4 := testInput{
		"NATS",
		false,
		false,
	}

	testCase := []struct {
		inp      testInput
		expected string
	}{
		{inp: ti1, expected: "^(?:(?!\\b(%s|%s)\\b).)*$"},
		{inp: ti2, expected: "^(?:(?!\\b(%s)\\b).)*$"},
		{inp: ti3, expected: ".*\\b(%s|%s)\\b.*"},
		{inp: ti4, expected: ".*\\b(%s)\\b.*"},
	}

	for _, tCase := range testCase {
		t.Run(tCase.inp.word, func(t *testing.T) {
			res := generateRegex(tCase.inp.word, tCase.inp.bothRegs, tCase.inp.revert)
			if res != tCase.expected {
				fmt.Printf("Wrong result %s for %v. Expected: %s", res, tCase.inp, tCase.expected)
			}
		})
	}
}
