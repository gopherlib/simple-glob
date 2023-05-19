package debug

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/gopherlib/simple-glob/match"
)

func Graphviz(pattern string, m match.Matcher) string {
	return fmt.Sprintf(`digraph G {graph[label="%s"];%s}`, pattern, graphvizInternal(m, fmt.Sprintf("%x", rand.Int63())))
}

func graphvizInternal(m match.Matcher, id string) string {
	buf := &bytes.Buffer{}

	switch matcher := m.(type) {
	case match.BTree:
		_, _ = fmt.Fprintf(buf, `"%s"[label="%s"];`, id, matcher.Value.String())
		for _, m := range []match.Matcher{matcher.Left, matcher.Right} {
			switch n := m.(type) {
			case nil:
				rnd := rand.Int63()
				_, _ = fmt.Fprintf(buf, `"%x"[label="<nil>"];`, rnd)
				_, _ = fmt.Fprintf(buf, `"%s"->"%x";`, id, rnd)

			default:
				sub := fmt.Sprintf("%x", rand.Int63())
				_, _ = fmt.Fprintf(buf, `"%s"->"%s";`, id, sub)
				_, _ = fmt.Fprintf(buf, graphvizInternal(n, sub))
			}
		}
	default:
		_, _ = fmt.Fprintf(buf, `"%s"[label="%s"];`, id, m.String())
	}

	return buf.String()
}
