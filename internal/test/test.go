package test

import "testing"

func AssertEqual(t *testing.T, actual, expected interface{}, message string) {
    if actual != expected {
        t.Errorf("Assertion failed: %s\nExpected: %v\nActual: %v", message, expected, actual)
    }
}

func AssertTrue(t *testing.T, actual bool, message string) {
    if !actual {
        t.Errorf("Assertion failed: %s\nExpected: true\nActual: false", message)
    }
}

func AssertFalse(t *testing.T, actual bool, message string) {
    if actual {
        t.Errorf("Assertion failed: %s\nExpected: false\nActual: true", message)
    }
}