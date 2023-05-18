package match

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Text raw represents raw string to match
type Text struct {
	Str         string
	RunesLength int
	BytesLength int
	Segments    []int
}

func NewText(s string) Text {
	return Text{
		Str:         s,
		RunesLength: utf8.RuneCountInString(s),
		BytesLength: len(s),
		Segments:    []int{len(s)},
	}
}

func (t Text) Match(s string) bool {
	return t.Str == s
}

func (t Text) Len() int {
	return t.RunesLength
}

func (t Text) Index(s string) (int, []int) {
	index := strings.Index(s, t.Str)
	if index == -1 {
		return -1, nil
	}

	return index, t.Segments
}

func (t Text) String() string {
	return fmt.Sprintf("<text:`%v`>", t.Str)
}
