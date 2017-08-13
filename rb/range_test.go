package rb

import (
	"testing"
	h "github.com/ArieShout/go-ruby/rb/testing"
)

func TestRange_First(t *testing.T) {
	r := NewRange(1, 3)
	if !h.SliceEquals(r.First(1), []int{1}) {
		t.Log("First(1) does not equal to []int{1}")
		t.Fail()
	}
	if !h.SliceEquals(r.First(5), []int{1, 2, 3}) {
		t.Log("First(5) does not equal to []int{1, 2, 3}")
		t.Fail()
	}
}

func TestRange_Each(t *testing.T) {
	r := make([]int, 0, 10)
	NewRange(1, 3).Each(func(i int) {
		if i == 3 {
			Break()
		}
		r = append(r, i)
	})
	h.AssertTrue(t, "Break on Each", h.SliceEquals(r, []int{1, 2}))
}

func TestRange_ToS(t *testing.T) {
	h.AssertTrue(t, "Empty range", NewRange(1, 0).ToS() == "[]")
	h.AssertTrue(t, "One element range", NewRange(1, 1).ToS() == "[1]")
	h.AssertTrue(t, "Two elements range", NewRange(1, 2).ToS() == "[1, 2]")
	h.AssertTrue(t, "Exclude end range", NewRangeExclusive(1, 3).ToS() == "[1, 2]")
}
