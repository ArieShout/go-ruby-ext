package rb

import (
	"regexp"
	"strconv"
	"reflect"
	"bytes"
)

type MatchData struct {
	str    string
	regexp *regexp.Regexp
	index  []int
}

func (m MatchData) checkIndex(n int) {
	if n < 0 || n*2+1 >= len(m.index) {
		panic("index " + strconv.Itoa(n) + " out of matches")
	}
}

func (m MatchData) Begin(n int) int {
	m.checkIndex(n)
	return m.index[n*2]
}

func (m MatchData) Captures() []*string {
	size := m.Size()
	caps := make([]*string, size)
	for i := 1; i <= size; i++ {
		begin, end := m.index[i*2], m.index[i*2+1]
		if begin >= 0 {
			caps[i-1] = &m.str[begin:end]
		}
	}
	return caps
}

func (m MatchData) End(n int) int {
	m.checkIndex(n)
	return m.index[n*2+1]
}

func (m MatchData) Group(n int) *string {
	m.checkIndex(n)
	begin, end := m.index[n*2], m.index[n*2+1]
	if begin >= 0 {
		return &m.str[begin:end]
	}
	return nil
}

func (m MatchData) IsEql(rhs MatchData) bool {
	return m.str == rhs.str && reflect.DeepEqual(m.regexp, rhs.regexp) && reflect.DeepEqual(m.index, rhs.index)
}

func safeString(str *string) string {
	if str == nil {
		return "nil"
	} else {
		return "\"" + strconv.Quote(*str) + "\""
	}
}

func (m MatchData) Inspect() string {
	var buf bytes.Buffer
	buf.WriteString(`#<MatchData `)
	buf.WriteString(safeString(m.Group(0)))

	names := m.regexp.SubexpNames()
	for i, size := 1, m.Size(); i <= size; i++ {
		buf.WriteRune(' ')
		name := names[i]
		if name != "" {
			buf.WriteString(name)
		}
		buf.WriteRune(':')
		buf.WriteString(safeString(m.Group(i)))
	}
	buf.WriteRune('>')

	return buf.String()
}

func (m MatchData) Length() int {
	return m.Size()
}

func (m MatchData) NamedCaptures() map[string]*string {
	caps := make(map[string]*string)
	for i, name := range m.regexp.SubexpNames() {
		if name != "" {
			caps[name] = m.Group(i)
		}
	}
	return caps
}

func (m MatchData) Names() []string {
	size := m.Size()
	names := make([]string, 0, size)
	source := m.regexp.SubexpNames()
	for i := 1; i <= size; i++ {
		name := source[i]
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

func (m MatchData) Offset(n int) []int {
	m.checkIndex(n)
	return []int{m.index[n*2], m.index[n*2+1]}
}

func (m MatchData) PreMatch() string {
	return m.str[:m.index[0]]
}

func (m MatchData) PostMatch() string {
	return m.str[m.index[1]+1:]
}

func (m MatchData) Regexp() *regexp.Regexp {
	return m.regexp
}

func (m MatchData) Size() int {
	indexLen := len(m.index)
	if indexLen > 1 {
		return indexLen/2 - 1
	} else {
		return 0
	}
}

func (m MatchData) OriginalString() string {
	return m.str
}

func (m MatchData) ToA() []*string {
	size := m.Size() + 1
	arr := make([]*string, size)
	for i := 0; i < size; i++ {
		begin, end := m.index[i*2], m.index[i*2+1]
		if begin >= 0 {
			arr[i-1] = &m.str[begin:end]
		}
	}
	return arr
}

func (m MatchData) String() string {
	group := m.Group(0)
	if group == nil {
		return ""
	}
	return *group
}

func (m MatchData) ValuesAt(index ...int) []*string {
	arr := make([]*string, len(index))
	for i := 0; i < len(index); i++ {
		arr[i] = m.Group(i)
	}
	return arr
}
