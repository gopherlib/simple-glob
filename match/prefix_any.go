package match

import (
	"fmt"
	"strings"
	"unicode/utf8"

	sutil "github.com/gopherlib/simple-glob/util/strings"
)

type PrefixAny struct {
	Prefix     string
	Separators []rune
}

func NewPrefixAny(s string, sep []rune) PrefixAny {
	return PrefixAny{s, sep}
}

func (a PrefixAny) Index(s string) (int, []int) {
	idx := strings.Index(s, a.Prefix)
	if idx == -1 {
		return -1, nil
	}

	n := len(a.Prefix)
	sub := s[idx+n:]
	i := sutil.IndexAnyRunes(sub, a.Separators)
	if i > -1 {
		sub = sub[:i]
	}

	seg := acquireSegments(len(sub) + 1)
	seg = append(seg, n)
	for i, r := range sub {
		seg = append(seg, n+i+utf8.RuneLen(r))
	}

	return idx, seg
}

func (a PrefixAny) Len() int {
	return lenNo
}

func (a PrefixAny) Match(s string) bool {
	if !strings.HasPrefix(s, a.Prefix) {
		return false
	}
	return sutil.IndexAnyRunes(s[len(a.Prefix):], a.Separators) == -1
}

func (a PrefixAny) String() string {
	return fmt.Sprintf("<prefix_any:%s![%s]>", a.Prefix, string(a.Separators))
}
