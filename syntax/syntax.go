package syntax

import (
	"github.com/gopher/simple-glob/syntax/ast"
	"github.com/gopher/simple-glob/syntax/lexer"
)

func Parse(s string) (*ast.Node, error) {
	return ast.Parse(lexer.NewLexer(s))
}

func Special(b byte) bool {
	return lexer.Special(b)
}
