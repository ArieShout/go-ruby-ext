package rb

import (
	"unicode/utf8"
)

type tr struct {
	pattern string
}

type trIter struct {
	tr           *tr
	buffer       [3]rune
	bufferIndex  int
	patternIndex int
	rng          Range
}

func (iter *trIter) fillBuffer() {
	for iter.bufferIndex < len(iter.buffer) && iter.patternIndex < len(iter.tr.pattern) {
		r, width := utf8.DecodeRuneInString(iter.tr.pattern[iter.patternIndex:])
		if r != utf8.RuneError {
			iter.buffer[iter.bufferIndex] = r
			iter.bufferIndex++
		}
		iter.patternIndex += width
	}
}

func (iter *trIter) drain(count int) {
	for i := count; i < len(iter.buffer); i++ {
		iter.buffer[i-count] = iter.buffer[i]
	}
	if iter.bufferIndex >= count {
		iter.bufferIndex -= count
	} else {
		iter.bufferIndex = 0
	}
}

func (iter *trIter) next() (r rune, ok bool) {
	if !iter.rng.IsEmpty() {
		r = rune(iter.rng.First())
		iter.rng = NewRange(iter.rng.First()+1, iter.rng.Last())
		return r, true
	}

	iter.fillBuffer()

	if iter.bufferIndex == 0 {
		return 0, false
	}

	r = iter.buffer[0]
	ok = true

	if iter.bufferIndex == 1 {
		iter.drain(1)
		return
	}

	if iter.bufferIndex == 2 {
		if r == '\\' {
			r = iter.buffer[1]
			iter.drain(2)
		} else {
			iter.drain(1)
		}
		return
	}

	if r == '\\' {
		r = iter.buffer[1]
		iter.drain(1)
		iter.fillBuffer()
	}

	if iter.bufferIndex < 3 || iter.buffer[1] != '-' {
		iter.drain(1)
		return
	}

	iter.rng = NewRange(int(iter.buffer[0]), int(iter.buffer[2]))
	iter.drain(3)
	return iter.next()
}

func (tr tr) IsNegative() bool {
	return len(tr.pattern) > 0 && tr.pattern[0] == '^'
}

func (tr tr) LazyChars() func() (rune, bool) {
	iter := trIter{&tr, [3]rune{}, 0, 0, NewRange(0, -1)}
	if tr.IsNegative() {
		iter.next()
	}
	return func() (rune, bool) {
		return iter.next()
	}
}

func (tr tr) Chars() []rune {
	chars := make([]rune, 0, 4)
	iter := tr.LazyChars()
	for {
		c, ok := iter()
		if !ok {
			break
		}
		chars = append(chars, c)
	}
	return chars
}
