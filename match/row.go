package match

import (
	"fmt"
)

type Row struct {
	Matchers    Matchers
	RunesLength int
	Segments    []int
}

func NewRow(len int, m ...Matcher) Row {
	return Row{
		Matchers:    Matchers(m),
		RunesLength: len,
		Segments:    []int{len},
	}
}

func (r Row) matchAll(s string) bool {
	var idx int
	for _, m := range r.Matchers {
		length := m.Len()

		var next, i int
		for next = range s[idx:] {
			i++
			if i == length {
				break
			}
		}

		if i < length || !m.Match(s[idx:idx+next+1]) {
			return false
		}

		idx += next + 1
	}

	return true
}

func (r Row) lenOk(s string) bool {
	var i int
	for range s {
		i++
		if i > r.RunesLength {
			return false
		}
	}
	return r.RunesLength == i
}

func (r Row) Match(s string) bool {
	return r.lenOk(s) && r.matchAll(s)
}

func (r Row) Len() (l int) {
	return r.RunesLength
}

func (r Row) Index(s string) (int, []int) {
	for i := range s {
		if len(s[i:]) < r.RunesLength {
			break
		}
		if r.matchAll(s[i:]) {
			return i, r.Segments
		}
	}
	return -1, nil
}

func (r Row) String() string {
	return fmt.Sprintf("<row_%d:[%s]>", r.RunesLength, r.Matchers)
}
