package ast

import (
	"errors"
	"fmt"

	"github.com/gopher/simple-glob/syntax/lexer"
)

type Lexer interface {
	Next() lexer.Token
}

type parseFn func(*Node, Lexer) (parseFn, *Node, error)

func Parse(lexer Lexer) (*Node, error) {
	var parser parseFn

	root := NewNode(KindPattern, nil)

	var (
		tree *Node
		err  error
	)
	for parser, tree = parserMain, root; parser != nil; {
		parser, tree, err = parser(tree, lexer)
		if err != nil {
			return nil, err
		}
	}

	return root, nil
}

func parserMain(tree *Node, lex Lexer) (parseFn, *Node, error) {
	for {
		token := lex.Next()
		switch token.Type {
		case lexer.EOF:
			return nil, tree, nil

		case lexer.Error:
			return nil, tree, errors.New(token.Raw)

		case lexer.Text:
			Insert(tree, NewNode(KindText, Text{token.Raw}))
			return parserMain, tree, nil

		case lexer.Any:
			Insert(tree, NewNode(KindAny, nil))
			return parserMain, tree, nil

		case lexer.Separator:
			p := NewNode(KindPattern, nil)
			Insert(tree.Parent, p)

			return parserMain, p, nil

		default:
			return nil, tree, fmt.Errorf("unexpected token: %s", token)
		}
	}

	//return nil, tree, fmt.Errorf("unknown error")
}
