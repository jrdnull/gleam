package gleam

import (
	"strings"
	"unicode/utf8"
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	}

	return t.val
}

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF

	tokenLeftParen
	tokenRightParen
	tokenNumber
	tokenSymbol
)

const eof = -1

type lexer struct {
	input      string
	start      int // start position of this token
	pos        int // current position in the input
	width      int // width of the last rune read from input
	parenDepth int
	tokens     []token
}

func lex(input string) []token {
	l := &lexer{
		input: input,
	}

	l.run()

	return l.tokens
}

// stateFn represents the current state and returns the next
type stateFn func(*lexer) stateFn

// run lexes the input executing stateFns til the state is nil
func (l *lexer) run() {
	for state := lexWhitespace; state != nil; {
		state = state(l)
	}
}

func (l *lexer) emit(t tokenType) {
	l.tokens = append(l.tokens, token{t, l.input[l.start:l.pos]})
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r

}

func (l *lexer) ignore() {
	l.start = l.pos
}

// acceptRun consumes a run of runes for the valid set
func (l *lexer) acceptRun(valid string) {
	for l.accept(valid) {
	}
}

// accept consumes the next rune if its in the valid set
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}

	l.backup()

	return false
}

type runePredicate func(rune) bool

func (l *lexer) acceptUntil(p runePredicate) {
	for !p(l.next()) {
	}

	l.backup()
}

// lexWhitespace is the initial stateFn, eating up whitespace
// until something more interesting comes along or we finish
func lexWhitespace(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(tokenEOF)

			return nil
		case isWhitespace(r):
			l.ignore()

			return lexWhitespace
		case r == ')':
			return lexRightParen
		case r == '(':
			return lexLeftParen
		case '0' <= r && r <= '9', r == '+' || r == '-':
			if (r == '+' || r == '-') && isWhitespace(l.peek()) {
				l.backup()

				return lexSymbol
			}

			l.backup()

			return lexNumber
		default:
			l.backup()

			return lexSymbol
		}
	}
}

func lexLeftParen(l *lexer) stateFn {
	l.parenDepth++
	l.emit(tokenLeftParen)

	return lexWhitespace
}

func lexRightParen(l *lexer) stateFn {
	l.parenDepth--
	l.emit(tokenRightParen)

	return lexWhitespace
}

const digits = "0123456789"

func lexNumber(l *lexer) stateFn {
	l.accept("+-") // Optional sign
	l.acceptRun(digits)

	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(tokenNumber)

	return lexWhitespace
}

func lexSymbol(l *lexer) stateFn {
	l.acceptUntil(func(r rune) bool {
		return isWhitespace(r) || r == ')' || r == eof
	})

	l.emit(tokenSymbol)

	return lexWhitespace
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}
