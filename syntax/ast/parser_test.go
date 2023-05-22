package ast

import (
	"reflect"
	"testing"

	"github.com/gopherlib/simple-glob/syntax/lexer"
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
				{Type: lexer.Text, Raw: "abc"},
				{Type: lexer.EOF, Raw: ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "abc"}),
			),
		},
		{
			//pattern: "a*c",
			testName: "a*c",
			tokens: []lexer.Token{
				{Type: lexer.Text, Raw: "a"},
				{Type: lexer.Any, Raw: "*"},
				{Type: lexer.Text, Raw: "c"},
				{Type: lexer.EOF, Raw: ""},
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
				{Type: lexer.Text, Raw: "a"},
				{Type: lexer.Any, Raw: "*"},
				{Type: lexer.Any, Raw: "*"},
				{Type: lexer.Text, Raw: "c"},
				{Type: lexer.EOF, Raw: ""},
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
				{Type: lexer.Text, Raw: "a?c"},
				{Type: lexer.EOF, Raw: ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "a?c"}),
			),
		},
		{
			//pattern: "[!a-z]",
			testName: "[!a-z]",
			tokens: []lexer.Token{
				{Type: lexer.Text, Raw: "[!a-z]"},
				{Type: lexer.EOF, Raw: ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "[!a-z]"}),
			),
		},
		{
			//pattern: "[az]",
			testName: "[az]",
			tokens: []lexer.Token{
				{Type: lexer.Text, Raw: "[az]"},
				{Type: lexer.EOF, Raw: ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "[az]"}),
			),
		},
		{
			//pattern: "{a,z}",
			testName: "{a,z}",
			tokens: []lexer.Token{
				{Type: lexer.Text, Raw: "{a,z}"},
				{Type: lexer.EOF, Raw: ""},
			},
			tree: NewNode(KindPattern, nil,
				NewNode(KindText, Text{Text: "{a,z}"}),
			),
		},
		{
			//pattern: "/{z,ab}*",
			testName: "/{z,ab}*",
			tokens: []lexer.Token{
				{Type: lexer.Text, Raw: "/{z,ab}"},
				{Type: lexer.Any, Raw: "*"},
				{Type: lexer.EOF, Raw: ""},
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
				{Type: lexer.Text, Raw: "{a,{x,y},?,[a-z],[!qwe]}"},
				{Type: lexer.EOF, Raw: ""},
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
