package match

import (
	"fmt"
	"unicode/utf8"
)

type BTree struct {
	Value            Matcher
	Left             Matcher
	Right            Matcher
	ValueLengthRunes int
	LeftLengthRunes  int
	RightLengthRunes int
	LengthRunes      int
}

func NewBTree(Value, Left, Right Matcher) (tree BTree) {
	tree.Value = Value
	tree.Left = Left
	tree.Right = Right

	lenOk := true
	if tree.ValueLengthRunes = Value.Len(); tree.ValueLengthRunes == -1 {
		lenOk = false
	}

	if Left != nil {
		if tree.LeftLengthRunes = Left.Len(); tree.LeftLengthRunes == -1 {
			lenOk = false
		}
	}

	if Right != nil {
		if tree.RightLengthRunes = Right.Len(); tree.RightLengthRunes == -1 {
			lenOk = false
		}
	}

	if lenOk {
		tree.LengthRunes = tree.LeftLengthRunes + tree.ValueLengthRunes + tree.RightLengthRunes
	} else {
		tree.LengthRunes = -1
	}

	return tree
}

func (t BTree) Len() int {
	return t.LengthRunes
}

// Index todo?
func (t BTree) Index(s string) (index int, segments []int) {
	//inputLen := len(s)
	//// try to cut unnecessary parts
	//// by knowledge of length of right and left part
	//offset, limit := t.offsetLimit(inputLen)
	//for offset < limit {
	//	// search for matching part in substring
	//	vi, segments := t.Value.Index(s[offset:limit])
	//	if index == -1 {
	//		return -1, nil
	//	}
	//	if t.Left == nil {
	//		if index != offset {
	//			return -1, nil
	//		}
	//	} else {
	//		left := s[:offset+vi]
	//		i := t.Left.IndexSuffix(left)
	//		if i == -1 {
	//			return -1, nil
	//		}
	//		index = i
	//	}
	//	if t.Right != nil {
	//		for _, seg := range segments {
	//			right := s[:offset+vi+seg]
	//		}
	//	}

	//	l := s[:offset+index]
	//	var left bool
	//	if t.Left != nil {
	//		left = t.Left.Index(l)
	//	} else {
	//		left = l == ""
	//	}
	//}

	return -1, nil
}

func (t BTree) Match(s string) bool {
	inputLen := len(s)
	// try to cut unnecessary parts
	// by knowledge of length of right and left part
	offset, limit := t.offsetLimit(inputLen)

	for offset < limit {
		// search for matching part in substring
		index, segments := t.Value.Index(s[offset:limit])
		if index == -1 {
			releaseSegments(segments)
			return false
		}

		l := s[:offset+index]
		var left bool
		if t.Left != nil {
			left = t.Left.Match(l)
		} else {
			left = l == ""
		}

		if left {
			for i := len(segments) - 1; i >= 0; i-- {
				length := segments[i]

				var right bool
				var r string
				// if there is no string for the right branch
				if inputLen <= offset+index+length {
					r = ""
				} else {
					r = s[offset+index+length:]
				}

				if t.Right != nil {
					right = t.Right.Match(r)
				} else {
					right = r == ""
				}

				if right {
					releaseSegments(segments)
					return true
				}
			}
		}

		_, step := utf8.DecodeRuneInString(s[offset+index:])
		offset += index + step

		releaseSegments(segments)
	}

	return false
}

func (t BTree) offsetLimit(inputLen int) (offset int, limit int) {
	// t.Length, t.RLen and t.LLen are values meaning the length of runes for each part
	// here we manipulating byte length for better optimizations
	// but these checks still works, cause minLen of 1-rune string is 1 byte.
	if t.LengthRunes != -1 && t.LengthRunes > inputLen {
		return 0, 0
	}
	if t.LeftLengthRunes >= 0 {
		offset = t.LeftLengthRunes
	}
	if t.RightLengthRunes >= 0 {
		limit = inputLen - t.RightLengthRunes
	} else {
		limit = inputLen
	}
	return offset, limit
}

func (t BTree) String() string {
	const n string = "<nil>"
	var l, r string
	if t.Left == nil {
		l = n
	} else {
		l = t.Left.String()
	}
	if t.Right == nil {
		r = n
	} else {
		r = t.Right.String()
	}

	return fmt.Sprintf("<btree:[%s<-%s->%s]>", l, t.Value, r)
}
