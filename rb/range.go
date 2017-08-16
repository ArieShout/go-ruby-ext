package rb

import (
	"bytes"
	"strconv"
	"fmt"
)

type Range struct {
	first, last int
	excludeEnd  bool
}

func NewRange(begin, end int) Range {
	return Range{begin, end, false}
}

func NewRangeExclusive(begin, end int) Range {
	return Range{begin, end, true}
}

func (r Range) OpEquals(obj interface{}) bool {
	return r.IsEql(obj)
}

func (r Range) OpCaseEquals(obj interface{}) bool {
	return r.IsCover(obj)
}

func (r Range) actualEnd() int {
	if r.excludeEnd {
		return r.last - 1
	}
	return r.last
}

func (r Range) Bsearch(pred func(int) bool) (ret int, found bool) {
	begin, end := r.first, r.actualEnd()
	if begin > end {
		return
	}
	if pred(begin) {
		ret = begin
		found = true
		return
	}
	if !pred(end) {
		return
	}
	for begin < end {
		mid := begin + (end-begin)/2
		if pred(mid) {
			end = mid
		} else {
			begin = mid + 1
		}
	}
	return begin, true
}

func (r Range) BsearchAny(pred func(int) int) (int, bool) {
	return r.Bsearch(func(i int) bool {
		return pred(i) == 0
	})
}

func (r Range) IsCover(obj interface{}) bool {
	if rhs, ok := obj.(int); ok {
		return rhs >= r.first && (rhs < r.last || (rhs == r.last && !r.excludeEnd))
	}
	return false
}

func (r Range) IsEmpty() bool {
	return r.Size() == 0
}

func (r Range) Each(action func(int)) {
	defer RecoverBreak("")
	for i, end := r.first, r.actualEnd(); i <= end; i++ {
		action(i)
	}
}

func (r Range) IsEql(obj interface{}) bool {
	if rhs, ok := obj.(Range); ok {
		return r == rhs
	}
	return false
}

func (r Range) ExcludeEnd() bool {
	return r.excludeEnd
}

func (r Range) First() int {
	return r.first
}

func (r Range) FirstSlice(limit int) []int {
	ret := make([]int, 0, limit)
	for count, i, end := 0, r.first, r.actualEnd(); count < limit && i <= end; i++ {
		ret = append(ret, i)
		count++
	}
	return ret
}

func (r Range) IsInclude(obj interface{}) bool {
	return r.IsCover(obj)
}

func (r Range) Inspect() string {
	if r.excludeEnd {
		return fmt.Sprintf("%d...%d", r.first, r.last)
	} else {
		return fmt.Sprintf("%d..%d", r.first, r.last)
	}
}

func (r Range) Last() int {
	return r.last
}

func (r Range) Max() (max int, ok bool) {
	last := r.actualEnd()
	return last, r.First() <= last
}

func (r Range) IsMember(obj interface{}) bool {
	return r.IsCover(obj)
}

func (r Range) Min() (min int, ok bool) {
	min = r.first
	ok = r.first <= r.actualEnd()
	return
}

func (r Range) Size() int {
	size := r.actualEnd() - r.first + 1
	if size < 0 {
		size = 0
	}
	return size
}

func (r Range) Step(step int, action func(int)) {
	defer RecoverBreak("")
	for i, end := r.first, r.actualEnd(); i < end; i += step {
		action(i)
	}
}

func (r Range) ToS() string {
	return r.String()
}

func (r Range) String() string {
	var buf bytes.Buffer
	buf.WriteRune('[')
	end := r.actualEnd()
	if r.first <= end {
		buf.WriteString(strconv.Itoa(r.first))
		for i := r.first + 1; i <= end; i++ {
			buf.WriteRune(',')
			buf.WriteRune(' ')
			buf.WriteString(strconv.Itoa(i))
		}
	}
	buf.WriteRune(']')
	return buf.String()
}
