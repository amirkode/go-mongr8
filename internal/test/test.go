package test

import "testing"

func AssertEqual(t *testing.T, actual, expected interface{}, message string) {
    if actual != expected {
        t.Errorf("Assertion failed (Equal): %s\nExpected: %v\nActual: %v", message, expected, actual)
    }
}

func AssertNotEqual(t *testing.T, actual, expected interface{}, message string) {
    if actual == expected {
        t.Errorf("Assertion failed (Not Equal): %s\nExpected: %v\nActual: %v", message, expected, actual)
    }
}

func AssertTrue(t *testing.T, actual bool, message string) {
    if !actual {
        t.Errorf("Assertion failed (True): %s\nExpected: true\nActual: false", message)
    }
}

func AssertFalse(t *testing.T, actual bool, message string) {
    if actual {
        t.Errorf("Assertion failed (False): %s\nExpected: false\nActual: true", message)
    }
}