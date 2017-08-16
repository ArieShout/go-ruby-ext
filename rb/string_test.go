package rb

import (
	"testing"
	"fmt"
	"regexp"
	"github.com/stretchr/testify/assert"
)

func TestString_String(t *testing.T) {
	base := "abc"
	str := NewString(base)
	var stringer fmt.Stringer = fmt.Stringer(str)
	assert.Equal(t, base, stringer.String(), "String() should return the underlying string")
}

func TestString_Concat(t *testing.T) {
	assert.Equal(t, "ab", NewString("a").Concat("b").Value, `"a".Concat("b")`)
	assert.Equal(t, "ab99", NewString("ab").Concat(99).Value, `"ab".Concat(99)`)
	assert.Equal(t, "ab99d", NewString("ab99").Concat(rune(100)).Value, `"ab99".Concat(rune(100))`)
	assert.Equal(t, "ab99de", NewString("ab99d").Concat(String{"e"}).Value, `"ab99d".Concat(String{"e"})`)
	assert.Equal(t, "ab99de1", NewString("ab99de").Concat(byte(1)).Value, `"ab99de".Concat(byte(1))`)
	assert.Equal(t, "atrue", NewString("a").Concat(true).Value, `"a".Concat(true)`)
	assert.Equal(t, "abc", NewString("").Concat("a", "b", "c").Value, `"".Concat("a", "b", "c")`)
}

func TestString_OpSubscript(t *testing.T) {
	str := NewString("abc红宝石")

	sub, ok := str.OpSubscript(3)
	assert.True(t, ok, "utf8 index OK")
	assert.Equal(t, "红", sub.Value, "utf8 index")

	sub, ok = str.OpSubscript(-1)
	assert.True(t, ok, "utf8 negative index OK")
	assert.Equal(t, "石", sub.Value, "utf8 negative index")

	sub, ok = str.OpSubscript2(2, 2)
	assert.True(t, ok, "start, length OK")
	assert.Equal(t, "c红", sub.Value, "start, length")

	sub, ok = str.OpSubscript(NewRange(1, 4))
	assert.True(t, ok, "range OK")
	assert.Equal(t, "bc红宝", sub.Value, "range")

	sub, ok = str.OpSubscript(regexp.MustCompile(`c.`))
	assert.True(t, ok, "regexp OK")
	assert.Equal(t, "c红", sub.Value, "regexp")

	sub, ok = str.OpSubscript2(regexp.MustCompile(`c(.)`), 1)
	assert.True(t, ok, "regexp, int OK")
	assert.Equal(t, "红", sub.Value, "regexp, int")

	sub, ok = str.OpSubscript(NewString("红"))
	assert.True(t, ok, "String OK")
	assert.Equal(t, "红", sub.Value, "String")

	sub, ok = str.OpSubscript("红")
	assert.True(t, ok, "string OK")
	assert.Equal(t, "红", sub.Value, "string")
}

func TestString_Center(t *testing.T) {
	str := NewString("hello")

	assert.Equal(t, "hello", str.Center(4).Value, `"hello".Center(4)`)
	assert.Equal(t, "hello ", str.Center(6).Value, `"hello".Center(6)`)
	assert.Equal(t, " hello ", str.Center(7).Value, `"hello".Center(7)`)
}

func TestString_CenterWith(t *testing.T) {
	assert.Equal(t, "1231231hello12312312", NewString("hello").Center2(20, NewString("123")).Value, `"hello".Center2(20, "123")`)
	assert.Equal(t, "我a我我我", NewString("a我").Center2(5, NewString("我")).Value, "Multibytes padding")
}

func TestString_Chomp(t *testing.T) {
	assert.Equal(t, "abc", NewString("abc").Chomp().Value, "No update if not ended with LF")
	assert.Equal(t, "abc", NewString("abc\n").Chomp().Value, "Remove trailing LF")
}

func TestString_Chomp1(t *testing.T) {
	assert.Equal(t, "a", NewString("abc").Chomp1(NewString("bc")).Value, "Remove trailing sequence")
}

func TestString_Count(t *testing.T) {
	a := NewString("hello world")
	assert.Equal(t, 5, a.Count(NewString("lo")), `"hello world".Count("lo")`)
	assert.Equal(t, 2, a.Count(NewString("lo"), NewString("o")), `"hello world".Count("lo", "o")`)
	assert.Equal(t, 4, a.Count(NewString("hello"), NewString("^l")), `"hello world".Count("hello", "^l")`)
	assert.Equal(t, 4, a.Count(NewString("ej-m")), `"hello world".Count("ej-m")`)
	assert.Equal(t, 4, NewString("hello^world").Count(NewString("\\^aeiou")), `"hello^world".Count("\\\\^aeiou")`)
	assert.Equal(t, 4, NewString("hello-world").Count(NewString("a\\-eo")), `"hello-world".Count("a\\\\-eo")`)

	a = NewString("hello world\\r\\n")
	assert.Equal(t, 2, a.Count(NewString("\\")), `"hello world".Count("\\\\")`)
	assert.Equal(t, 0, a.Count(NewString("\\A")), `"hello world".Count("\\\\A")`)
	assert.Equal(t, 3, a.Count(NewString("X-\\w")), `"hello world".Count("X-\\\\w")`)
}

func TestString_Delete(t *testing.T) {
	assert.Equal(t, "heo", NewString("hello").Delete(NewString("l"), NewString("lo")).Value, `"hello".Delete("l", "lo")`)
	assert.Equal(t, "he", NewString("hello").Delete(NewString("lo")).Value, `"hello".Delete("lo")`)
	assert.Equal(t, "hell", NewString("hello").Delete(NewString("aeiou"), NewString("^e")).Value, `"hello".Delete("aeiou", "^e")`)
	assert.Equal(t, "ho", NewString("hello").Delete(NewString("ej-m")).Value, `"hello".Delete("ej-m")`)
}

func TestString_Dump(t *testing.T) {
	assert.Equal(t, `"hello \n ''"`, NewString("hello \n ''").Dump().Value, `\"hello \\n ''\".Dump()`)
	assert.Equal(t, `"a\v"`, NewString("a\013").Dump().Value, `\"hello \\n ''\".Dump()`)
	assert.Equal(t, `"a\xC4"`, NewString("a\xC4").Dump().Value, `\"hello \\n ''\".Dump()`)
}
