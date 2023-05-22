package compiler

import (
	"fmt"

	"github.com/gopherlib/simple-glob/match"
	"github.com/gopherlib/simple-glob/syntax/ast"
	"github.com/gopherlib/simple-glob/util/runes"
)

func optimizeMatcher(matcher match.Matcher) match.Matcher {
	switch m := matcher.(type) {

	case match.Any:
		if len(m.Separators) == 0 {
			return m
		}

	case match.BTree:
		m.Left = optimizeMatcher(m.Left)
		m.Right = optimizeMatcher(m.Right)

		r, ok := m.Value.(match.Text)
		if !ok {
			return m
		}

		var (
			leftNil  = m.Left == nil
			rightNil = m.Right == nil
		)
		if leftNil && rightNil {
			return match.NewText(r.Str)
		}

		la, leftAny := m.Left.(match.Any)
		ra, rightAny := m.Right.(match.Any)

		switch {
		case rightNil && leftAny:
			return match.NewSuffixAny(r.Str, la.Separators)

		case leftNil && rightAny:
			return match.NewPrefixAny(r.Str, ra.Separators)
		}

		return m
	}

	return matcher
}

func compileMatchers(matchers []match.Matcher) (match.Matcher, error) {
	if len(matchers) == 0 {
		return nil, fmt.Errorf("compile error: need at least one matcher")
	}
	if len(matchers) == 1 {
		return matchers[0], nil
	}
	if m := glueMatchers(matchers); m != nil {
		return m, nil
	}

	idx := -1
	maxLen := -1
	var val match.Matcher
	for i, matcher := range matchers {
		if l := matcher.Len(); l != -1 && l >= maxLen {
			maxLen = l
			idx = i
			val = matcher
		}
	}

	if val == nil { // not found matcher with static length
		r, err := compileMatchers(matchers[1:])
		if err != nil {
			return nil, err
		}
		return match.NewBTree(matchers[0], nil, r), nil
	}

	left := matchers[:idx]
	var right []match.Matcher
	if len(matchers) > idx+1 {
		right = matchers[idx+1:]
	}

	var l, r match.Matcher
	var err error
	if len(left) > 0 {
		l, err = compileMatchers(left)
		if err != nil {
			return nil, err
		}
	}

	if len(right) > 0 {
		r, err = compileMatchers(right)
		if err != nil {
			return nil, err
		}
	}

	return match.NewBTree(val, l, r), nil
}

func glueMatchers(matchers []match.Matcher) match.Matcher {
	if m := glueMatchersAsEvery(matchers); m != nil {
		return m
	}
	if m := glueMatchersAsRow(matchers); m != nil {
		return m
	}
	return nil
}

func glueMatchersAsRow(matchers []match.Matcher) match.Matcher {
	if len(matchers) <= 1 {
		return nil
	}

	var (
		c []match.Matcher
		l int
	)
	for _, matcher := range matchers {
		if ml := matcher.Len(); ml == -1 {
			return nil
		} else {
			c = append(c, matcher)
			l += ml
		}
	}
	return match.NewRow(l, c...)
}

func glueMatchersAsEvery(matchers []match.Matcher) match.Matcher {
	if len(matchers) <= 1 {
		return nil
	}

	var (
		hasAny    bool
		separator []rune
	)

	for i, matcher := range matchers {
		var sep []rune

		switch m := matcher.(type) {
		case match.Any:
			sep = m.Separators
			hasAny = true

		default:
			return nil
		}

		// initialize
		if i == 0 {
			separator = sep
		}

		if runes.Equal(sep, separator) {
			continue
		}

		return nil
	}

	if hasAny {
		return match.NewAny(separator)
	}

	return nil
}

func minimizeMatchers(matchers []match.Matcher) []match.Matcher {
	var done match.Matcher
	var left, right, count int

	for l := 0; l < len(matchers); l++ {
		for r := len(matchers); r > l; r-- {
			if glued := glueMatchers(matchers[l:r]); glued != nil {
				var swap bool

				if done == nil {
					swap = true
				} else {
					cl, gl := done.Len(), glued.Len()
					swap = cl > -1 && gl > -1 && gl > cl
					swap = swap || count < r-l
				}

				if swap {
					done = glued
					left = l
					right = r
					count = r - l
				}
			}
		}
	}

	if done == nil {
		return matchers
	}

	next := append(append([]match.Matcher{}, matchers[:left]...), done)
	if right < len(matchers) {
		next = append(next, matchers[right:]...)
	}

	if len(next) == len(matchers) {
		return next
	}

	return minimizeMatchers(next)
}

func commonChildren(nodes []*ast.Node) (commonLeft, commonRight []*ast.Node) {
	if len(nodes) <= 1 {
		return
	}

	// find node that has least number of children
	idx := leastChildren(nodes)
	if idx == -1 {
		return
	}
	tree := nodes[idx]
	treeLength := len(tree.Children)

	// allocate max able size for rightCommon slice
	// to get ability insert elements in reverse order (from end to start)
	// without sorting
	commonRight = make([]*ast.Node, treeLength)
	lastRight := treeLength // will use this to get results as commonRight[lastRight:]

	var (
		breakLeft   bool
		breakRight  bool
		commonTotal int
	)
	for i, j := 0, treeLength-1; commonTotal < treeLength && j >= 0 && !(breakLeft && breakRight); i, j = i+1, j-1 {
		treeLeft := tree.Children[i]
		treeRight := tree.Children[j]

		for k := 0; k < len(nodes) && !(breakLeft && breakRight); k++ {
			// skip least children node
			if k == idx {
				continue
			}

			restLeft := nodes[k].Children[i]
			restRight := nodes[k].Children[j+len(nodes[k].Children)-treeLength]

			breakLeft = breakLeft || !treeLeft.Equal(restLeft)

			// disable searching for right common parts, if left part is already overlapping
			breakRight = breakRight || (!breakLeft && j <= i)
			breakRight = breakRight || !treeRight.Equal(restRight)
		}

		if !breakLeft {
			commonTotal++
			commonLeft = append(commonLeft, treeLeft)
		}
		if !breakRight {
			commonTotal++
			lastRight = j
			commonRight[j] = treeRight
		}
	}

	commonRight = commonRight[lastRight:]

	return
}

func leastChildren(nodes []*ast.Node) int {
	min := -1
	idx := -1
	for i, n := range nodes {
		if idx == -1 || (len(n.Children) < min) {
			min = len(n.Children)
			idx = i
		}
	}
	return idx
}

func compileTreeChildren(tree *ast.Node, sep []rune) ([]match.Matcher, error) {
	var matchers []match.Matcher
	for _, desc := range tree.Children {
		m, err := compile(desc, sep)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, optimizeMatcher(m))
	}
	return matchers, nil
}

func compile(tree *ast.Node, sep []rune) (m match.Matcher, err error) {
	switch tree.Kind {

	case ast.KindPattern:
		if len(tree.Children) == 0 {
			return match.NewNothing(), nil
		}
		matchers, err := compileTreeChildren(tree, sep)
		if err != nil {
			return nil, err
		}
		m, err = compileMatchers(minimizeMatchers(matchers))
		if err != nil {
			return nil, err
		}

	case ast.KindAny:
		m = match.NewAny(sep)

	case ast.KindNothing:
		m = match.NewNothing()

	case ast.KindText:
		t := tree.Value.(ast.Text)
		m = match.NewText(t.Text)

	default:
		return nil, fmt.Errorf("could not compile tree: unknown node type")
	}

	return optimizeMatcher(m), nil
}

func Compile(tree *ast.Node, sep []rune) (match.Matcher, error) {
	m, err := compile(tree, sep)
	if err != nil {
		return nil, err
	}

	return m, nil
}
