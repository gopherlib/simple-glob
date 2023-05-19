package glob

import (
	"github.com/gopherlib/simple-glob/compiler"
	"github.com/gopherlib/simple-glob/syntax"
)

// Glob represents compiled glob pattern.
type Glob interface {
	Match(string) bool
}

// Compile creates Glob for given pattern and strings (if any present after pattern) as separators.
// The pattern syntax is:
//
//	pattern:
//	    { term }
//
//	term:
//	    `*`         matches any sequence of non-separator characters
func Compile(pattern string, separators ...rune) (Glob, error) {
	ast, err := syntax.Parse(pattern)
	if err != nil {
		return nil, err
	}

	matcher, err := compiler.Compile(ast, separators)
	if err != nil {
		return nil, err
	}

	return matcher, nil
}

// MustCompile is the same as Compile, except that if Compile returns error, this will panic
func MustCompile(pattern string, separators ...rune) Glob {
	g, err := Compile(pattern, separators...)
	if err != nil {
		panic(err)
	}

	return g
}
