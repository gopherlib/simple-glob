package match

import (
	"fmt"

	"github.com/gopher/simple-glob/util/strings"
)

type Any struct {
	Separators []rune
}

func NewAny(s []rune) Any {
	return Any{s}
}

func (a Any) Match(s string) bool {
	return strings.IndexAnyRunes(s, a.Separators) == -1
}

func (a Any) Index(s string) (int, []int) {
	found := strings.IndexAnyRunes(s, a.Separators)
	switch found {
	case -1:
	case 0:
		return 0, segments0
	default:
		s = s[:found]
	}

	segments := acquireSegments(len(s))
	for i := range s {
		segments = append(segments, i)
	}
	segments = append(segments, len(s))

	return 0, segments
}

func (a Any) Len() int {
	return lenNo
}

func (a Any) String() string {
	return fmt.Sprintf("<any:![%s]>", string(a.Separators))
}
