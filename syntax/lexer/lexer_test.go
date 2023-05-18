package lexer

import (
	"testing"
)

func TestLexGood(t *testing.T) {
	for id, test := range []struct {
		pattern string
		items   []Token
	}{
		{
			pattern: "",
			items: []Token{
				{EOF, ""},
			},
		},
		{
			pattern: `\*`,
			items: []Token{
				{Text, `\`},
				{Any, "*"},
				{EOF, ""},
			},
		},
		{
			pattern: "hello",
			items: []Token{
				{Text, "hello"},
				{EOF, ""},
			},
		},
		{
			pattern: "/{rate,[0-9]]}*",
			items: []Token{
				{Text, "/{rate,[0-9]]}"},
				{Any, "*"},
				{EOF, ""},
			},
		},
		{
			pattern: "hello,world",
			items: []Token{
				{Text, "hello,world"},
				{EOF, ""},
			},
		},
		{
			pattern: "hello\\,world",
			items: []Token{
				{Text, "hello\\,world"},
				{EOF, ""},
			},
		},
		{
			pattern: "hello\\{world",
			items: []Token{
				{Text, "hello\\{world"},
				{EOF, ""},
			},
		},
		{
			pattern: "hello?",
			items: []Token{
				{Text, "hello?"},
				{EOF, ""},
			},
		},
		{
			pattern: "hellof*",
			items: []Token{
				{Text, "hellof"},
				{Any, "*"},
				{EOF, ""},
			},
		},
		{
			pattern: "hello**",
			items: []Token{
				{Text, "hello"},
				{Any, "*"},
				{Any, "*"},
				{EOF, ""},
			},
		},
		{
			pattern: "[日-語]",
			items: []Token{
				{Text, "[日-語]"},
				{EOF, ""},
			},
		},
		{
			pattern: "[!日-語]",
			items: []Token{
				{Text, "[!日-語]"},
				{EOF, ""},
			},
		},
		{
			pattern: "[日本語]",
			items: []Token{
				{Text, "[日本語]"},
				{EOF, ""},
			},
		},
		{
			pattern: "[!日本語]",
			items: []Token{
				{Text, "[!日本語]"},
				{EOF, ""},
			},
		},
		{
			pattern: "{a,b}",
			items: []Token{
				{Text, "{a,b}"},
				{EOF, ""},
			},
		},
		{
			pattern: "/{z,ab}*",
			items: []Token{
				{Text, "/{z,ab}"},
				{Any, "*"},
				{EOF, ""},
			},
		},
		{
			pattern: "{[!日-語],*,?,{a,b,\\c}}",
			items: []Token{
				{Text, "{[!日-語],"},
				{Any, "*"},
				{Text, ",?,{a,b,\\c}}"},
				{EOF, ""},
			},
		},
	} {
		t.Run(test.pattern, func(t *testing.T) {
			lexer := NewLexer(test.pattern)
			for i, exp := range test.items {
				act := lexer.Next()
				if act.Type != exp.Type {
					t.Errorf("#%d %q: wrong %d-th item type: exp: %q; act: %q\n\t(%s vs %s)", id, test.pattern, i, exp.Type, act.Type, exp, act)
				}
				if act.Raw != exp.Raw {
					t.Errorf("#%d %q: wrong %d-th item contents: exp: %q; act: %q\n\t(%s vs %s)", id, test.pattern, i, exp.Raw, act.Raw, exp, act)
				}
			}
		})
	}
}
