package testing

import (
	"testing"
	"fmt"
)

func SliceEquals(a, b []int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range(a) {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func AssertTrue(t *testing.T, msg string, pred bool) {
	if !pred {
		t.Log(msg)
		t.Fail()
	}
}

func FailWithMessage(t *testing.T, message string) {
	t.Log(message)
	t.Fail()
}

func AssertEquals(t *testing.T, msg string, expected string, actual string) {
	if expected != actual {
		t.Log(fmt.Sprintf("%s, expected %s, got %s", msg, expected, actual))
		t.Fail()
	}
}
