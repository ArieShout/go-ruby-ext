package rb

import (
	"fmt"
	"bytes"
	"strings"
	"regexp"
	"unicode/utf8"
	"strconv"
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

func (str String) OpLtLt(args ...interface{}) String {
	return str.Concat(args)
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
		strLen := str.Length()
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
		first := rng.First()
		last := rng.Last()
		if last < 0 {
			last = str.Length() + last
		}
		if rng.ExcludeEnd() {
			last -= 1
		}
		if last < first {
			return NewString(""), true
		}
		buf := make([]rune, 0, last-first+1)
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
			return str.OpSubscript(NewRangeExclusive(start, start+length))
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

func fillToLength(str String, length int) string {
	var buf bytes.Buffer
	strLen := str.Length()
	for ; length > 0; length -= strLen {
		if length > strLen {
			buf.WriteString(str.Value)
		} else {
			substr, _ := str.OpSubscript2(0, length)
			buf.WriteString(substr.Value)
		}
	}
	return buf.String()
}

func (str String) Center(width int) String {
	return str.Center2(width, NewString(" "))
}

func (str String) Center2(width int, padstr String) String {
	if padstr.IsEmpty() {
		panic("Empty padding")
	}
	strLen := str.Length()
	leftPad := (width - strLen) / 2
	rightPad := width - strLen - leftPad
	return NewString(fillToLength(padstr, leftPad) + str.Value + fillToLength(padstr, rightPad))
}

func (str String) Chars() []String {
	chars := make([]String, str.Length())
	i := 0
	for _, c := range str.Value {
		chars[i] = NewString(string(c))
		i++
	}
	return chars
}

func (str String) Chars_() []rune {
	chars := make([]rune, str.Length())
	i := 0
	for _, c := range str.Value {
		chars[i] = c
		i++
	}
	return chars
}

func (str String) Chomp() String {
	return str.Chomp1(NewString("\n"))
}

func (str String) Chomp1(separator String) String {
	if separator.IsEmpty() || !str.EndWith(separator) {
		return str
	}
	result, _ := str.OpSubscript(NewRangeExclusive(0, -separator.Length()))
	return result
}

func trSetupTable(charSet String, includes, excludes map[rune]bool, intersect *bool) {
	tr := tr{charSet.Value}
	target := includes
	isIntersect := *intersect
	if tr.IsNegative() {
		target = excludes
		isIntersect = false
	} else {
		if isIntersect {
			for k, v := range target {
				if v {
					target[k] = false
				} else {
					delete(target, k)
				}
			}
		}
		*intersect = true
	}
	iter := tr.LazyChars()
	for {
		r, ok := iter()
		if !ok {
			break
		}
		if isIntersect {
			_, ok = target[r]
			if ok {
				target[r] = true
			}
		} else {
			target[r] = true
		}
	}
}

func (str String) Count(charset String, otherCharsets ...String) int {
	includes := make(map[rune]bool)
	excludes := make(map[rune]bool)

	intersect := false

	trSetupTable(charset, includes, excludes, &intersect)
	for _, s := range otherCharsets {
		trSetupTable(s, includes, excludes, &intersect)
	}

	count := 0
	for _, r := range str.Value {
		if _, ok := excludes[r]; ok {
			continue
		}

		if v, ok := includes[r]; ok && v {
			count++
		}
	}
	return count
}

// TODO crypt

func (str String) Delete(charset String, otherCharset ...String) String {
	includes := make(map[rune]bool)
	excludes := make(map[rune]bool)
	intersect := false

	trSetupTable(charset, includes, excludes, &intersect)
	for _, s := range otherCharset {
		trSetupTable(s, includes, excludes, &intersect)
	}

	remain := make([]rune, 0, str.Length())
	for _, r := range str.Value {
		if _, ok := excludes[r]; !ok {
			if v, ok := includes[r]; ok && v {
				continue
			}
		}

		remain = append(remain, r)
	}

	return NewString(string(remain))
}

func (str String) DeletePrefix(prefix String) String {
	if str.StartWith(prefix) {
		return NewString(str.Value[len(prefix.Value):])
	}
	return NewString(str.Value)
}

func (str String) Downcase() String {
	return NewString(strings.ToLower(str.Value))
}

func (str String) Dump() String {
	width := 0
	var r rune
	bufLen := 2 // ""
	for i := 0; i < len(str.Value); i += width {
		r, width = utf8.DecodeRuneInString(str.Value[i:])
		if r == utf8.RuneError {
			if width == 0 {
				break
			}
			if width == 1 {
				bufLen += 4
				continue
			}
		}
		switch r {
		case '"', '\\', '\n', '\r', '\t', '\f', '\013', '\010', '\007', '\033':
			bufLen += 2
		default:
			if r <= 0x7F {
				if strconv.IsPrint(r) {
					bufLen += 1
				} else {
					bufLen += 4 // \xNN
				}
			} else if r <= 0xFFFF {
				bufLen += 6 // \uXXXX
			} else {
				bufLen += 10 // \uXXXXXXXX
			}
		}
	}

	buffer := bytes.NewBuffer(make([]byte, 0, bufLen))
	buffer.WriteByte('"')
	for i := 0; i < len(str.Value); i += width {
		r, width = utf8.DecodeRuneInString(str.Value[i:])
		if r == utf8.RuneError {
			if width == 0 {
				break
			} else if width == 1 {
				buffer.WriteString(fmt.Sprintf("\\x%02X", str.Value[i]))
				continue
			}
		}
		switch r {
		case '"', '\\':
			buffer.WriteByte('\\')
			buffer.WriteByte(byte(r))
		case '\n':
			buffer.WriteByte('\\')
			buffer.WriteByte('n')
		case '\r':
			buffer.WriteByte('\\')
			buffer.WriteByte('r')
		case '\t':
			buffer.WriteByte('\\')
			buffer.WriteByte('t')
		case '\f':
			buffer.WriteByte('\\')
			buffer.WriteByte('f')
		case '\013':
			buffer.WriteByte('\\')
			buffer.WriteByte('v')
		case '\010':
			buffer.WriteByte('\\')
			buffer.WriteByte('b')
		case '\007':
			buffer.WriteByte('\\')
			buffer.WriteByte('a')
		case '\033':
			buffer.WriteByte('\\')
			buffer.WriteByte('e')
		default:
			if r <= 0x7F && strconv.IsPrint(r) {
				buffer.WriteByte(byte(r))
			} else {
				buffer.WriteByte('\\')
				if r <= 0x7F {
					buffer.WriteString(fmt.Sprintf("x%02X", r))
				} else if r <= 0xFFFF {
					buffer.WriteString(fmt.Sprintf("u%04X", r))
				} else {
					buffer.WriteString(fmt.Sprintf("u%08X", r))
				}
			}
		}
	}
	buffer.WriteByte('"')
	return NewString(buffer.String())
}

func (str String) EachByte(action func(byte)) (ret String) {
	ret = str
	defer RecoverBreak("")
	for i := 0; i < len(str.Value); i++ {
		action(str.Value[i])
	}
	return
}

func (str String) EachChar(action func(string)) (ret String) {
	return str.EachCodepoint(func(r rune) {
		action(string(r))
	})
}

func (str String) EachCodepoint(action func(rune)) (ret String) {
	ret = str
	defer RecoverBreak("")
	var r rune
	width := 0
	for i := 0; i < len(str.Value); i += width {
		r, width = utf8.DecodeRuneInString(str.Value[i:])
		if r == utf8.RuneError && width <= 1 {
			panic("Invalid char at position " + strconv.Itoa(i))
		}
		action(r)
	}
	return
}

func (str String) EachLine(separator String, action func(String)) (ret String) {
	ret = str
	if separator.Value == "" {
		separator = NewString("\n")
	}
	for i := 0; i < len(str.Value); {
		idx := strings.Index(str.Value[i:], separator.Value)
		if idx >= 0 {
			action(NewString(str.Value[i:i+idx]))
			i += idx + len(separator.Value)
			if i == len(str.Value) {
				// the separator is at the end of the String
				action(NewString(""))
			}
		} else {
			action(NewString(str.Value[i:]))
			break
		}
	}
	return
}

// TODO encode?

func (str String) EndWith(suffix String, otherSuffixes ...String) bool {
	if strings.HasSuffix(str.Value, suffix.Value) {
		return true
	}
	for _, s := range otherSuffixes {
		if strings.HasSuffix(str.Value, s.Value) {
			return true
		}
	}
	return false
}

func (str String) GetByte(index int) byte {
	return str.Value[index]
}

type subValues interface {
	getValue()
}

func (str String) Gsub(re regexp.Regexp, replacement interface{}) String {
	indexes := re.FindAllStringSubmatchIndex(str.Value, -1)
	if indexes == nil {
		return NewString(str.Value)
	}
	return str
}

func (str String) IsEql(obj interface{}) bool {
	if rhs, ok := obj.(String); ok {
		return str.Value == rhs.Value
	}
	return false
}

func (str String) IsEmpty() bool {
	return str.Value == ""
}



func (str String) Length() int {
	return utf8.RuneCountInString(str.Value)
}

func (str String) Lines(separator String) []String {
	lines := make([]String, 0, 1)
	str.EachLine(separator, func(line String) {
		lines = append(lines, line)
	})
	return lines
}

func (str String) StartWith(prefix String, otherPrefixes ...String) bool {
	if strings.HasPrefix(str.Value, prefix.Value) {
		return true
	}
	for _, s := range otherPrefixes {
		if strings.HasPrefix(str.Value, s.Value) {
			return true
		}
	}
	return false
}

func (str String) Upcase() String {
	return NewString(strings.ToUpper(str.Value))
}
