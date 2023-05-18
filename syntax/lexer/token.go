package lexer

import "fmt"

type TokenType int

const (
	EOF TokenType = iota
	Error
	Text
	Any
	Separator
)

func (tt TokenType) String() string {
	switch tt {
	case EOF:
		return "eof"

	case Error:
		return "error"

	case Text:
		return "text"

	case Any:
		return "any"

	case Separator:
		return "separator"

	default:
		return "undef"
	}
}

type Token struct {
	Type TokenType
	Raw  string
}

func (t Token) String() string {
	return fmt.Sprintf("%v<%q>", t.Type, t.Raw)
}
