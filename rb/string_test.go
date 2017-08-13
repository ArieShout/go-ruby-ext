package rb

import (
	"testing"
	"fmt"
	h "github.com/ArieShout/go-ruby-ext/rb/testing"
	"regexp"
)

func TestString_String(t *testing.T) {
	base := "abc"
	str := NewString(base)
	var stringer fmt.Stringer = fmt.Stringer(str)
	if stringer.String() != base {
		t.Log("String() should return the underlying string")
		t.Fail()
	}
}

func TestString_Concat(t *testing.T) {
	str := NewString("a")
	str = str.Concat("b")
	if str.Value != "ab" {
		t.Log("Concat string")
		t.Fail()
	}
	str = str.Concat(99)
	if str.Value != "abc" {
		t.Log("Concat int code point")
		t.Fail()
	}
	str = str.Concat(rune(100))
	if str.Value != "abcd" {
		t.Log("Concat rune code point")
		t.Fail()
	}
	str = str.Concat(String{"e"})
	if str.Value != "abcde" {
		t.Log("Concat fmt.Stringer")
		t.Fail()
	}
	str = str.Concat(byte(1))
	if str.Value != "abcde1" {
		t.Log("Concat other numeric")
		t.Fail()
	}
	str = str.Concat(true)
	if str.Value != "abcde1true" {
		t.Log("Concat bool")
		t.Fail()
	}
	str = String{}.Concat("a", "b", "c")
	if str.Value != "abc" {
		t.Log("Concat multiple")
		t.Fail()
	}
}

func TestString_OpSubscript(t *testing.T) {
	str := NewString("abc红宝石")

	sub, ok := str.OpSubscript(3)
	h.AssertTrue(t, "utf8 index OK", ok)
	h.AssertEquals(t,"utf8 index", "红", sub.Value)

	sub, ok = str.OpSubscript(-1)
	h.AssertTrue(t, "utf8 negative index OK", ok)
	h.AssertEquals(t, "utf8 negative index", "石", sub.Value)

	sub, ok = str.OpSubscript2(2, 2)
	h.AssertTrue(t, "start, length OK", ok)
	h.AssertEquals(t, "start, length", "c红", sub.Value)

	sub, ok = str.OpSubscript(NewRange(1, 4))
	h.AssertTrue(t, "range OK", ok)
	h.AssertEquals(t, "range", "bc红宝", sub.Value)

	sub, ok = str.OpSubscript(regexp.MustCompile(`c.`))
	h.AssertTrue(t, "regexp OK", ok)
	h.AssertEquals(t, "regexp", "c红", sub.Value)

	sub, ok = str.OpSubscript2(regexp.MustCompile(`c(.)`), 1)
	h.AssertTrue(t, "regexp, int OK", ok)
	h.AssertEquals(t, "regexp, int", "红", sub.Value)

	sub, ok = str.OpSubscript(NewString("红"))
	h.AssertTrue(t, "String OK", ok)
	h.AssertEquals(t, "String", "红", sub.Value)

	sub, ok = str.OpSubscript("红")
	h.AssertTrue(t, "string OK", ok)
	h.AssertEquals(t, "string", "红", sub.Value)
}