package testing

import "testing"

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
