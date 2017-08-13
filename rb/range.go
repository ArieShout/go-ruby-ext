package rb

import (
	"bytes"
	"strconv"
	"fmt"
)

type Range struct {
	Begin, End int
	ExcludeEnd bool
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
	if r.ExcludeEnd {
		return r.End - 1
	}
	return r.End
}

func (r Range) Bsearch(pred func(int) bool) (ret int, found bool) {
	begin, end := r.Begin, r.actualEnd()
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
		return rhs >= r.Begin && (rhs < r.End || (rhs == r.End && !r.ExcludeEnd))
	}
	return false
}

func (r Range) Each(action func(int)) {
	defer RecoverBreak("")
	for i, end := r.Begin, r.actualEnd(); i <= end; i++ {
		action(i)
	}
}

func (r Range) IsEql(obj interface{}) bool {
	if rhs, ok := obj.(Range); ok {
		return r == rhs
	}
	return false
}

func (r Range) First(limit int) []int {
	ret := make([]int, 0, limit)
	for count, i, end := 0, r.Begin, r.actualEnd(); count < limit && i <= end; i++ {
		ret = append(ret, i)
		count++
	}
	return ret
}

func (r Range) FirstElem() (ret int, ok bool) {
	end := r.actualEnd()
	if r.Begin > end {
		return
	}
	return r.Begin, true
}

func (r Range) IsInclude(obj interface{}) bool {
	return r.IsCover(obj)
}

func (r Range) Inspect() string {
	if r.ExcludeEnd {
		return fmt.Sprintf("%d...%d", r.Begin, r.End)
	} else {
		return fmt.Sprintf("%d..%d", r.Begin, r.End)
	}
}

func (r Range) Last() (end int, ok bool) {
	end = r.actualEnd()
	ok = r.Begin <= end
	return
}

func (r Range) Max() (max int, ok bool) {
	return r.Last()
}

func (r Range) IsMember(obj interface{}) bool {
	return r.IsCover(obj)
}

func (r Range) Min() (min int, ok bool) {
	min = r.Begin
	ok = r.Begin <= r.actualEnd()
	return
}

func (r Range) Size() int {
	size := r.actualEnd() - r.Begin + 1
	if size < 0 {
		size = 0
	}
	return size
}

func (r Range) Step(step int, action func(int)) {
	defer RecoverBreak("")
	for i, end := r.Begin, r.actualEnd(); i < end; i += step {
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
	if r.Begin <= end {
		buf.WriteString(strconv.Itoa(r.Begin))
		for i := r.Begin + 1; i <= end; i++ {
			buf.WriteRune(',')
			buf.WriteRune(' ')
			buf.WriteString(strconv.Itoa(i))
		}
	}
	buf.WriteRune(']')
	return buf.String()
}