package gleam

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	}

	return i.val
}

type itemType int

const (
	itemError itemType = iota
	itemEOF

	itemLeftParen
	itemRightParen
	itemNumber
	itemSymbol
)

const eof = -1

type lexer struct {
	input      string
	start      int // start position of this item
	pos        int // current position in the input
	width      int // width of the last rune read from input
	parenDepth int
	items      chan item
}

func lex(input string) (*lexer, chan item) {
	l := &lexer{
		input: input,
		items: make(chan item),
	}

	go l.run()

	return l, l.items
}

// stateFn represents the current state and returns the next
type stateFn func(*lexer) stateFn

// run lexes the input executing stateFns til the state is nil
func (l *lexer) run() {
	for state := lexWhitespace; state != nil; {
		state = state(l)
	}

	close(l.items) // no more tokens will be sent
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
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

// errorf emits an error token and returns a nil stateFn to end lexing
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}

	return nil
}

// lexWhitespace is the initial stateFn, eating up whitespace
// until something more interesting comes along or we finish
func lexWhitespace(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			if l.parenDepth != 0 {
				return l.errorf("unclosed left paren")
			}

			l.emit(itemEOF)

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
	l.emit(itemLeftParen)

	return lexWhitespace
}

func lexRightParen(l *lexer) stateFn {
	if l.parenDepth == 0 {
		return l.errorf("unexpected right paren")
	}

	l.parenDepth--
	l.emit(itemRightParen)

	return lexWhitespace
}

const digits = "0123456789"

func lexNumber(l *lexer) stateFn {
	l.accept("+-") // Optional sign
	l.acceptRun(digits)

	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(itemNumber)

	return lexWhitespace
}

func lexSymbol(l *lexer) stateFn {
	l.acceptUntil(func(r rune) bool {
		return isWhitespace(r) || r == eof
	})

	l.emit(itemSymbol)

	return lexWhitespace
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}
