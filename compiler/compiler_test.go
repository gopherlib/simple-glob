package compiler

import (
	"reflect"
	"testing"

	"github.com/gopherlib/simple-glob/match"
	"github.com/gopherlib/simple-glob/match/debug"
	"github.com/gopherlib/simple-glob/syntax/ast"
)

var separators = []rune{'.'}

func TestCommonChildren(t *testing.T) {
	for i, test := range []struct {
		nodes []*ast.Node
		left  []*ast.Node
		right []*ast.Node
	}{
		{
			nodes: []*ast.Node{
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"z"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
				),
			},
		},
		{
			nodes: []*ast.Node{
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"z"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
				),
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"b"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
				),
			},
			left: []*ast.Node{
				ast.NewNode(ast.KindText, ast.Text{"a"}),
			},
			right: []*ast.Node{
				ast.NewNode(ast.KindText, ast.Text{"c"}),
			},
		},
		{
			nodes: []*ast.Node{
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"b"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
					ast.NewNode(ast.KindText, ast.Text{"d"}),
				),
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"b"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
					ast.NewNode(ast.KindText, ast.Text{"d"}),
				),
			},
			left: []*ast.Node{
				ast.NewNode(ast.KindText, ast.Text{"a"}),
				ast.NewNode(ast.KindText, ast.Text{"b"}),
			},
			right: []*ast.Node{
				ast.NewNode(ast.KindText, ast.Text{"c"}),
				ast.NewNode(ast.KindText, ast.Text{"d"}),
			},
		},
		{
			nodes: []*ast.Node{
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"b"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
				),
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"b"}),
					ast.NewNode(ast.KindText, ast.Text{"b"}),
					ast.NewNode(ast.KindText, ast.Text{"c"}),
				),
			},
			left: []*ast.Node{
				ast.NewNode(ast.KindText, ast.Text{"a"}),
				ast.NewNode(ast.KindText, ast.Text{"b"}),
			},
			right: []*ast.Node{
				ast.NewNode(ast.KindText, ast.Text{"c"}),
			},
		},
		{
			nodes: []*ast.Node{
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"d"}),
				),
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"d"}),
				),
				ast.NewNode(ast.KindNothing, nil,
					ast.NewNode(ast.KindText, ast.Text{"a"}),
					ast.NewNode(ast.KindText, ast.Text{"e"}),
				),
			},
			left: []*ast.Node{
				ast.NewNode(ast.KindText, ast.Text{"a"}),
			},
			right: []*ast.Node{},
		},
	} {
		left, right := commonChildren(test.nodes)
		if !nodesEqual(left, test.left) {
			t.Errorf("[%d] left, right := commonChildren(); left = %v; want %v", i, left, test.left)
		}
		if !nodesEqual(right, test.right) {
			t.Errorf("[%d] left, right := commonChildren(); right = %v; want %v", i, right, test.right)
		}
	}
}

func nodesEqual(a, b []*ast.Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i, av := range a {
		if !av.Equal(b[i]) {
			return false
		}
	}
	return true
}

func TestGlueMatchers(t *testing.T) {
	for id, test := range []struct {
		testName string
		in       []match.Matcher
		exp      match.Matcher
	}{
		{
			"any",
			[]match.Matcher{
				match.NewAny(separators),
			},
			match.NewAny(separators),
		},
		{
			"a",
			[]match.Matcher{
				match.NewAny([]rune{'a'}),
			},
			match.NewAny([]rune{'a'}),
		},
	} {
		act, err := compileMatchers(test.in)
		if err != nil {
			t.Errorf("#%d convert matchers error: %s", id, err)
			continue
		}

		if !reflect.DeepEqual(act, test.exp) {
			t.Errorf("#%d unexpected convert matchers result:\nact: %#v;\nexp: %#v", id, act, test.exp)
			continue
		}
	}
}

func TestCompileMatchers(t *testing.T) {
	for id, test := range []struct {
		in  []match.Matcher
		exp match.Matcher
	}{
		{
			[]match.Matcher{
				match.NewText("c"),
			},
			match.NewText("c"),
		},
		{
			[]match.Matcher{
				match.NewAny(nil),
				match.NewText("c"),
				match.NewAny(nil),
			},
			match.NewBTree(
				match.NewText("c"),
				match.NewAny(nil),
				match.NewAny(nil),
			),
		},
		{
			[]match.Matcher{
				match.NewText("c"),
			},
			match.NewText("c"),
		},
	} {
		act, err := compileMatchers(test.in)
		if err != nil {
			t.Errorf("#%d convert matchers error: %s", id, err)
			continue
		}

		if !reflect.DeepEqual(act, test.exp) {
			t.Errorf("#%d unexpected convert matchers result:\nact: %#v\nexp: %#v", id, act, test.exp)
			continue
		}
	}
}

func TestConvertMatchers(t *testing.T) {
	for id, test := range []struct {
		in, exp []match.Matcher
	}{
		{
			[]match.Matcher{
				match.NewText("c"),
				match.NewAny(nil),
			},
			[]match.Matcher{
				match.NewText("c"),
				match.NewAny(nil),
			},
		},
		{
			[]match.Matcher{
				match.NewText("c"),
				match.NewAny(nil),
				match.NewAny(nil),
			},
			[]match.Matcher{
				match.NewText("c"),
				match.NewAny(nil),
			},
		},
	} {
		act := minimizeMatchers(test.in)
		if !reflect.DeepEqual(act, test.exp) {
			t.Errorf("#%d unexpected convert matchers 2 result:\nact: %#v\nexp: %#v", id, act, test.exp)
			continue
		}
	}
}

func TestCompiler(t *testing.T) {
	for id, test := range []struct {
		testName string
		ast      *ast.Node
		result   match.Matcher
		sep      []rune
	}{
		{
			testName: "abc",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindText, ast.Text{"abc"}),
			),
			result: match.NewText("abc"),
		},
		{
			testName: "separators",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindAny, nil),
			),
			sep:    separators,
			result: match.NewAny(separators),
		},
		{
			testName: "any",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindAny, nil),
			),
			result: match.NewAny(nil),
		},
		{
			testName: "any_abc",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindAny, nil),
				ast.NewNode(ast.KindText, ast.Text{"abc"}),
			),
			sep:    separators,
			result: match.NewSuffixAny("abc", separators),
		},
		{
			testName: "abc_any",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindText, ast.Text{"abc"}),
				ast.NewNode(ast.KindAny, nil),
			),
			result: match.NewPrefixAny("abc", nil),
		},
		{
			testName: "abc_any_def",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindText, ast.Text{"abc"}),
				ast.NewNode(ast.KindAny, nil),
				ast.NewNode(ast.KindText, ast.Text{"def"}),
			),
			result: match.NewBTree(
				match.NewText("def"),
				match.NewPrefixAny("abc", nil),
				nil,
			),
		},
		{
			testName: "3any_abc_2any",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindAny, nil),
				ast.NewNode(ast.KindAny, nil),
				ast.NewNode(ast.KindAny, nil),
				ast.NewNode(ast.KindText, ast.Text{"abc"}),
				ast.NewNode(ast.KindAny, nil),
				ast.NewNode(ast.KindAny, nil),
			),
			sep: separators,
			result: match.NewBTree(
				match.NewText("abc"),
				match.NewAny(separators),
				match.NewAny(separators),
			),
		},
		{
			testName: "abc3",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindText, ast.Text{"abc"}),
			),
			result: match.NewText("abc"),
		},
		{
			testName: "p_abc",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindPattern, nil,
					ast.NewNode(ast.KindPattern, nil,
						ast.NewNode(ast.KindText, ast.Text{"abc"}),
					),
				),
			),
			result: match.NewText("abc"),
		},
		{
			testName: "4_abc",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindPattern, nil,
					ast.NewNode(ast.KindText, ast.Text{"abc"}),
				),
				ast.NewNode(ast.KindPattern, nil,
					ast.NewNode(ast.KindText, ast.Text{"abc"}),
				),
				ast.NewNode(ast.KindPattern, nil,
					ast.NewNode(ast.KindText, ast.Text{"abc"}),
				),
				ast.NewNode(ast.KindPattern, nil,
					ast.NewNode(ast.KindText, ast.Text{"abc"}),
				),
			),
			result: match.NewRow(
				12,
				match.NewText("abc"),
				match.NewText("abc"),
				match.NewText("abc"),
				match.NewText("abc"),
			),
		},
		{
			testName: "any_nil",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindAny, nil),
			),
			result: match.NewAny(nil),
		},
		{
			testName: "ghi_abc_ghi",
			ast: ast.NewNode(ast.KindPattern, nil,
				ast.NewNode(ast.KindPattern, nil,
					ast.NewNode(ast.KindText, ast.Text{"ghi"}),
				),
				ast.NewNode(ast.KindPattern, nil,
					ast.NewNode(ast.KindText, ast.Text{"abc"}),
					ast.NewNode(ast.KindText, ast.Text{"ghi"}),
				),
			),
			result: match.NewRow(
				9,
				match.NewText("ghi"),
				match.NewRow(
					6,
					match.NewText("abc"),
					match.NewText("ghi"),
				),
			),
		},
	} {
		t.Run(test.testName, func(t *testing.T) {
			m, err := Compile(test.ast, test.sep)
			if err != nil {
				t.Errorf("compilation error: %s", err)
			}

			if !reflect.DeepEqual(m, test.result) {
				t.Errorf("[%d] Compile():\nexp: %#v\nact: %#v\n\ngraphviz:\nexp:\n%s\nact:\n%s\n", id, test.result, m, debug.Graphviz("", test.result.(match.Matcher)), debug.Graphviz("", m.(match.Matcher)))
			}
		})
	}
}
