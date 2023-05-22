package match

import (
	"fmt"
	"strings"

	sutil "github.com/gopherlib/simple-glob/util/strings"
)

type SuffixAny struct {
	Suffix     string
	Separators []rune
}

func NewSuffixAny(s string, sep []rune) SuffixAny {
	return SuffixAny{s, sep}
}

func (a SuffixAny) Index(s string) (int, []int) {
	idx := strings.Index(s, a.Suffix)
	if idx == -1 {
		return -1, nil
	}

	i := sutil.LastIndexAnyRunes(s[:idx], a.Separators) + 1

	return i, []int{idx + len(a.Suffix) - i}
}

func (a SuffixAny) Len() int {
	return lenNo
}

func (a SuffixAny) Match(s string) bool {
	if !strings.HasSuffix(s, a.Suffix) {
		return false
	}
	return sutil.IndexAnyRunes(s[:len(s)-len(a.Suffix)], a.Separators) == -1
}

func (a SuffixAny) String() string {
	return fmt.Sprintf("<suffix_any:![%s]%s>", string(a.Separators), a.Suffix)
}
