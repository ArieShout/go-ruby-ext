package rb

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTr_Chars(t *testing.T) {
	assert.Equal(t, []rune{}, tr{""}.Chars(), "<empty>")
	assert.Equal(t, []rune{'a'}, tr{"a"}.Chars(), "a")
	assert.Equal(t, []rune{'a', 'b'}, tr{"ab"}.Chars(), "ab")
	assert.Equal(t, []rune{'a', 'b', 'c', 'd'}, tr{"ab-d"}.Chars(), "ab-d")
	assert.Equal(t, []rune{'\\'}, tr{"\\"}.Chars(), "\\\\")
	assert.Equal(t, []rune{'a', 'b', '-', 'd'}, tr{"ab\\-d"}.Chars(), "ab\\\\-d")
	assert.Equal(t, []rune{'a', 'b'}, tr{"a\\b"}.Chars(), "a\\\\b")
	assert.Equal(t, []rune{'a', '\\', 'b'}, tr{"a\\\\b"}.Chars(), "a\\\\\\\\b")
	assert.Equal(t, []rune{'a', '-', '.', '/'}, tr{"a\\--/"}.Chars(), "a\\\\--/")
	assert.Equal(t, []rune{'a', 'Z', '[', '\\'}, tr{"aZ-\\"}.Chars(), "aZ-\\\\")
	assert.Equal(t, []rune{'a', 'Z', '[', '\\', 'b'}, tr{"aZ-\\\\b"}.Chars(), "aZ-\\\\\\\\b")
	assert.Equal(t, []rune{'Z', '[', '\\', 'w'}, tr{"Z-\\w"}.Chars(), "Z-\\\\w")
}
