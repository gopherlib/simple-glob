package ast

import (
	"reflect"
	"testing"

	"github.com/gopher/simple-glob/syntax/lexer"
)

type stubLexer struct {
	tokens []lexer.Token
	pos    int
}

func (s *stubLexer) Next() (ret lexer.Token) {
	if s.pos == len(s.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	ret = s.tokens[s.pos]
	s.pos++
	return
}

func TestParseString(t *testing.T) {
	for id, test := range []struct {
		testName string
		tokens   []lexer.Token
		tree     *Node
	}{
		{
			//pattern: "abc",
			testName: "abc",
			tokens: []lexer.Token{
				{lexer.Text, "abc"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "abc"}),
			),
		},
		{
			//pattern: "a*c",
			testName: "a*c",
			tokens: []lexer.Token{
				{lexer.Text, "a"},
				{lexer.Any, "*"},
				{lexer.Text, "c"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "a"}),
				NewNode(KindAny, nil),
				NewNode(KindText, Text{Text: "c"}),
			),
		},
		{
			//pattern: "a**c",
			testName: "a**c",
			tokens: []lexer.Token{
				{lexer.Text, "a"},
				{lexer.Any, "*"},
				{lexer.Any, "*"},
				{lexer.Text, "c"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "a"}),
				NewNode(KindAny, nil),
				NewNode(KindAny, nil),
				NewNode(KindText, Text{Text: "c"}),
			),
		},
		{
			//pattern: "a?c",
			testName: "a?c",
			tokens: []lexer.Token{
				{lexer.Text, "a?c"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "a?c"}),
			),
		},
		{
			//pattern: "[!a-z]",
			testName: "[!a-z]",
			tokens: []lexer.Token{
				{lexer.Text, "[!a-z]"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "[!a-z]"}),
			),
		},
		{
			//pattern: "[az]",
			testName: "[az]",
			tokens: []lexer.Token{
				{lexer.Text, "[az]"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "[az]"}),
			),
		},
		{
			//pattern: "{a,z}",
			testName: "{a,z}",
			tokens: []lexer.Token{
				{lexer.Text, "{a,z}"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "{a,z}"}),
			),
		},
		{
			//pattern: "/{z,ab}*",
			testName: "/{z,ab}*",
			tokens: []lexer.Token{
				{lexer.Text, "/{z,ab}"},
				{lexer.Any, "*"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "/{z,ab}"}),
				NewNode(KindAny, nil),
			),
		},
		{
			//pattern: "{a,{x,y},?,[a-z],[!qwe]}",
			testName: "{a,{x,y},?,[a-z],[!qwe]}",
			tokens: []lexer.Token{
				{lexer.Text, "{a,{x,y},?,[a-z],[!qwe]}"},
				{lexer.EOF, ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "{a,{x,y},?,[a-z],[!qwe]}"}),
			),
		},
	} {
		t.Run(test.testName, func(t *testing.T) {
			l := &stubLexer{tokens: test.tokens}
			result, err := Parse(l)
			if err != nil {
				t.Errorf("[%d] unexpected error: %s", id, err)
			}
			if !reflect.DeepEqual(test.tree, result) {
				t.Errorf("[%d] Parse():\nact:\t%s\nexp:\t%s\n", id, result, test.tree)
			}
		})
	}
}
