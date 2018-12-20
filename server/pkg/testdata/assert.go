package testdata

import (
	"testing"
)

// AssertTrue asserts that isTrue is true
func AssertTrue(t *testing.T, isTrue bool) {
	t.Helper()

	if !isTrue {
		t.Fatalf("\nWant: true\n Got: %t\n", isTrue)
	}
}

// AssertNil asserts that err is nil
func AssertNil(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("\nWant: nil\n Got: %s\n", err)
	}
}

// AssertErr asserts that err is not nil
func AssertErr(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("\nWant: err\n Got: nil\n")
	}
}
