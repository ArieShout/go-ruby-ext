package rb

import (
	"fmt"
	"bytes"
	"strings"
	"regexp"
	"unicode/utf8"
)

type AsString interface {
	ToStr() String
}

type String struct {
	Value string
}

func NewString(str string) String {
	return String{str}
}

func (str String) String() string {
	return str.Value
}

func TryConvert(obj interface{}) (out String, ok bool) {
	ok = true
	switch obj.(type) {
	case string:
		out = String{obj.(string)}
	case String:
		out = obj.(String)
	default:
		out = NewString("")
		ok = false
	}
	return
}

func (str String) Format(params ...interface{}) String {
	return NewString(fmt.Sprintf(str.Value, params...))
}

func (str String) OpMultiply(times int) String {
	var buf bytes.Buffer
	for i := 0; i < times; i++ {
		buf.WriteString(str.Value)
	}
	return NewString(buf.String())
}

func (str String) OpAdd(rhs String) String {
	return String{str.Value + rhs.Value}
}

func (str String) Concat(args ...interface{}) String {
	var buf bytes.Buffer
	buf.WriteString(str.Value)
	for _, arg := range args {
		var value string
		switch arg.(type) {
		case rune:
			value = string(arg.(rune))
		case string:
			value = arg.(string)
		case fmt.Stringer:
			value = arg.(fmt.Stringer).String()
		default:
			value = fmt.Sprintf("%v", arg)
		}
		buf.WriteString(value)
	}
	return NewString(buf.String())
}

func (str String) OpSpaceShip(rhs String) int {
	return strings.Compare(str.Value, rhs.Value)
}

func (str String) OpEquals(obj interface{}) bool {
	if rhs, ok := obj.(String); ok {
		return str.Value == rhs.Value
	} else if rhs, ok := obj.(AsString); ok {
		return str.Value == rhs.ToStr().Value
	}
	return false
}

func (str String) OpCaseEquals(obj interface{}) bool {
	return str.OpEquals(obj)
}

func (str String) OpMatch(re regexp.Regexp) int {
	pos := re.FindStringIndex(str.Value)
	if pos == nil {
		return -1
	} else {
		return pos[0]
	}
}

func (str String) OpSubscript(arg interface{}) (ret String, found bool) {
	if index, ok := arg.(int); ok {
		strLen := utf8.RuneCountInString(str.Value)
		if index < 0 {
			// negative index
			index += strLen
		}
		if index < 0 || index >= strLen {
			return
		}
		i := 0
		for _, r := range str.Value {
			if i == index {
				return NewString(string(r)), true
			}
			i++
		}
		return
	}
	if rng, ok := arg.(Range); ok {
		buf := make([]rune, 0, rng.Size())
		first, ok := rng.First()
		if !ok {
			return
		}
		last, _ := rng.Last()
		i := 0
		for _, r := range str.Value {
			if i > last {
				break
			}
			if i >= first {
				buf = append(buf, r)
			}
			i++
		}
		ret = NewString(string(buf))
		return ret, true
	}
	if re, ok := arg.(*regexp.Regexp); ok {
		pos := re.FindStringIndex(str.Value)
		if pos == nil {
			return
		}
		return NewString(str.Value[pos[0]:pos[1]]), true
	}
	if s, ok := arg.(String); ok {
		if index := strings.Index(str.Value, s.Value); index >= 0 {
			return s, true
		}
		return
	}
	if s, ok := arg.(string); ok {
		if index := strings.Index(str.Value, s); index >= 0 {
			return NewString(s), true
		}
		return
	}
	panic("Argument type must be one of: int, Range, *Regexp, string")
}

func (str String) OpSubscript2(arg1, arg2 interface{}) (ret String, found bool) {
	if start, ok := arg1.(int); ok {
		if length, ok := arg2.(int); ok {
			if length < 0 {
				return
			}
			return str.OpSubscript(NewRangeExclusive(start, start + length))
		}
		goto TYPE_ERR
	}
	if re, ok := arg1.(*regexp.Regexp); ok {
		if capture, ok := arg2.(int); ok {
			if matches := re.FindStringSubmatch(str.Value); matches != nil {
				if capture < 0 {
					capture += len(matches)
				}
				if capture < 0 || capture >= len(matches) {
					return
				}
				return NewString(matches[capture]), true
			}
			return
		}
		goto TYPE_ERR
	}

	TYPE_ERR:
	panic("Arguments type must be one of: (int, int), (*Regexp, int)")
}

// TODO #[]=

func (str String) IsAsciiOnly() bool {
	for _, r := range str.Value {
		if r > 127 {
			return false
		}
	}
	return true
}

// TODO #b

func (str String) Bytes() []byte {
	return []byte(str.Value)
}

func (str String) Bytesize() int {
	return len(str.Value)
}

// TODO #byteslice

func (str String) Capitalize() String {
	if str.Value == "" {
		return str
	}
	_, width := utf8.DecodeRuneInString(str.Value)
	return NewString(strings.ToUpper(str.Value[0:width]) + strings.ToLower(str.Value[width:]))
}

func (str String) Casecmp(rhs String) int {
	return str.Downcase().OpSpaceShip(rhs.Downcase())
}

func (str String) IsCasecmp(rhs String) bool {
	return str.Downcase().OpEquals(rhs.Downcase())
}

func (str String) Center(width int, padstr string) String {

}

func (str String) Downcase() String {
	return NewString(strings.ToLower(str.Value))
}

func (str String) Upcase() String {
	return NewString(strings.ToUpper(str.Value))
}