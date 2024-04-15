/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package test

import (
	"fmt"
	"testing"
)

// Basic assertion with no testing scope
func Assert(expectsTrue bool, context string, unexpectedMessage string) {
	if !expectsTrue {
		panic(fmt.Sprintf("%s: %s", context, unexpectedMessage))
	}
}

func AssertEqual(t *testing.T, actual, expected interface{}, message string) {
	if actual != expected {
		msg := fmt.Sprintf("Assertion failed (Equal): %s\nExpected: %v\nActual: %v", message, expected, actual)
		t.Errorf(msg)
		panic(msg)
	}
}

func AssertNotEqual(t *testing.T, actual, expected interface{}, message string) {
	if actual == expected {
		msg := fmt.Sprintf("Assertion failed (Not Equal): %s\nExpected: %v\nActual: %v", message, expected, actual)
		t.Errorf(msg)
		panic(msg)
	}
}

func AssertTrue(t *testing.T, actual bool, message string) {
	if !actual {
		msg := fmt.Sprintf("Assertion failed (True): %s\nExpected: true\nActual: false", message)
		t.Errorf(msg)
		panic(msg)
	}
}

func AssertFalse(t *testing.T, actual bool, message string) {
	if actual {
		msg := fmt.Sprintf("Assertion failed (False): %s\nExpected: false\nActual: true", message)
		t.Errorf(msg)
		panic(msg)
	}
}
