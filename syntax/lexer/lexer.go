package lexer

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/gopherlib/simple-glob/util/runes"
)

const (
	charAny = '*'
)

var specials = []byte{
	charAny,
}

func Special(c byte) bool {
	return bytes.IndexByte(specials, c) != -1
}

type tokens []Token

func (i *tokens) shift() (ret Token) {
	ret = (*i)[0]
	copy(*i, (*i)[1:])
	*i = (*i)[:len(*i)-1]
	return
}

func (i *tokens) push(v Token) {
	*i = append(*i, v)
}

func (i *tokens) empty() bool {
	return len(*i) == 0
}

var eof rune = 0

type lexer struct {
	data string
	pos  int
	err  error

	tokens     tokens
	termsLevel int

	lastRune     rune
	lastRuneSize int
	hasRune      bool
}

func NewLexer(source string) *lexer {
	l := &lexer{
		data:   source,
		tokens: tokens(make([]Token, 0, 4)),
	}
	return l
}

func (l *lexer) Next() Token {
	if l.err != nil {
		return Token{Error, l.err.Error()}
	}
	if !l.tokens.empty() {
		return l.tokens.shift()
	}

	l.fetchItem()
	return l.Next()
}

func (l *lexer) peek() (r rune, w int) {
	if l.pos == len(l.data) {
		return eof, 0
	}

	r, w = utf8.DecodeRuneInString(l.data[l.pos:])
	if r == utf8.RuneError {
		l.errorf("could not read rune")
		r = eof
		w = 0
	}

	return
}

func (l *lexer) read() rune {
	if l.hasRune {
		l.hasRune = false
		l.seek(l.lastRuneSize)
		return l.lastRune
	}

	r, s := l.peek()
	l.seek(s)

	l.lastRune = r
	l.lastRuneSize = s

	return r
}

func (l *lexer) seek(w int) {
	l.pos += w
}

func (l *lexer) unread() {
	if l.hasRune {
		l.errorf("could not unread rune")
		return
	}
	l.seek(-l.lastRuneSize)
	l.hasRune = true
}

func (l *lexer) errorf(f string, v ...interface{}) {
	l.err = fmt.Errorf(f, v...)
}

func (l *lexer) inTerms() bool {
	return l.termsLevel > 0
}

func (l *lexer) termsEnter() {
	l.termsLevel++
}

func (l *lexer) termsLeave() {
	l.termsLevel--
}

var inTextBreakers = []rune{charAny}

func (l *lexer) fetchItem() {
	r := l.read()
	switch {
	case r == eof:
		l.tokens.push(Token{EOF, ""})
	case r == charAny:
		l.tokens.push(Token{Any, string(r)})
	default:
		l.unread()

		l.fetchText(inTextBreakers)
	}
}

func (l *lexer) fetchText(breakers []rune) {
	var data []rune

reading:
	for {
		r := l.read()
		if r == eof {
			break
		}

		if runes.IndexRune(breakers, r) != -1 {
			l.unread()
			break reading
		}

		data = append(data, r)
	}

	if len(data) > 0 {
		l.tokens.push(Token{Text, string(data)})
	}
}
