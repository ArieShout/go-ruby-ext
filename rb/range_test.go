package rb

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRange_First(t *testing.T) {
	r := NewRange(1, 3)
	assert.Equal(t, []int{1}, r.FirstSlice(1), `(1..3).FirstSlice(1)`)
	assert.Equal(t, []int{1, 2, 3}, r.FirstSlice(5), `(1..3).FirstSlice(5)`)
}

func TestRange_Each(t *testing.T) {
	r := make([]int, 0, 10)
	NewRange(1, 3).Each(func(i int) {
		if i == 3 {
			Break()
		}
		r = append(r, i)
	})
	assert.Equal(t, []int{1, 2}, r, "Break on Each")
}

func TestRange_ToS(t *testing.T) {
	assert.Equal(t, "[]", NewRange(1, 0).ToS(), "Empty range")
	assert.Equal(t, "[1]", NewRange(1, 1).ToS(), "One element range")
	assert.Equal(t, "[1, 2]", NewRange(1, 2).ToS(), "Two elements range")
	assert.Equal(t, "[1, 2]", NewRangeExclusive(1, 3).ToS(), "Exclude end range")
}
