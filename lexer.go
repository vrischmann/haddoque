package haddoque

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

//go:generate stringer -type=token
type token int

const (
	tokError token = iota
	tokEOF
	tokWhitespace

	tokField      // alphanumeric identifier starting with .
	tokIdentifier // alphanumeric identifier not starting with . (unused for now)

	// literals
	tokLiteralsBegin
	tokBool   // boolean constant
	tokChar   // character constant
	tokString // quoted string
	tokNumber // simple number
	tokLiteralsEnd

	// misc
	tokLparen   // (
	tokRparen   // )
	tokLbracket // [
	tokRbracket // ]
	tokComma    // ,

	// keywords
	tokKeywordsBegin
	tokWhere
	tokAnd
	tokOr
	tokIn
	tokContains
	tokKeywordsEnd

	// operators
	tokOperatorsBegin
	tokLt  // <
	tokLte // <=
	tokGt  // >
	tokGte // >=
	tokEq  // ==
	tokNeq // !=
	tokNot // !
	tokOperatorsEnd
)

// this is heavily based on the package text/template from the Go distribution

type lexeme struct {
	tok token
	pos int
	val string
}

func (l lexeme) String() string {
	return fmt.Sprintf("{%s, %d, %s}", l.tok, l.pos, l.val)
}

type lexStateFn func(*lexer) lexStateFn

type lexer struct {
	err   error
	input string
	start int
	pos   int
	width int
	state lexStateFn
	items chan lexeme
}

func newLexer(s string) *lexer {
	return &lexer{
		input: s,
		items: make(chan lexeme),
	}
}

const (
	eof = rune(-1)
)

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	return r
}

func (l *lexer) peek() rune {
	ch := l.next()
	l.backup()

	return ch
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) emit(tok token) {
	l.items <- lexeme{tok, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) lexStateFn {
	l.items <- lexeme{tokError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) lex() {
	go l.run()
}

func (l *lexer) nextLexeme() lexeme {
	return <-l.items
}

func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
}

func (l *lexer) atTerminator() bool {
	ch := l.peek()
	if isWhitespace(ch) {
		return true
	}

	switch ch {
	case eof, ',':
		return true
	}

	return false
}

func lexText(l *lexer) lexStateFn {
	switch ch := l.next(); {
	case ch == eof:
		l.emit(tokEOF)
		return nil
	case ch == ',':
		l.emit(tokComma)
	case ch == '.':
		return lexField
	case ch == '=':
		return lexEq
	case ch == '!':
		return lexNeq
	case ch == '(':
		l.emit(tokLparen)
	case ch == ')':
		l.emit(tokRparen)
	case ch == '[':
		l.emit(tokLbracket)
	case ch == ']':
		l.emit(tokRbracket)
	case ch == '<':
		return lexLt
	case ch == '>':
		return lexGt
	case ch == '+', ch == '-', ('0' <= ch && ch <= '9'):
		return lexNumber
	case ch == '"':
		return lexString
	case isAlphaNumeric(ch):
		return lexIdentifier
	case isWhitespace(ch):
		return lexSpace
	default:
		return l.errorf("query is malformed")
	}

	return lexText
}

func lexField(l *lexer) lexStateFn {
	var ch rune
	for {
		ch = l.next()
		if !isAlphaNumeric(ch) {
			l.backup()
			break
		}
	}

	l.emit(tokField)

	return lexText
}

func lexEq(l *lexer) lexStateFn {
	ch := l.next()
	if ch != '=' {
		return l.errorf("expected = after =")
	}

	l.emit(tokEq)

	return lexText
}

func lexNeq(l *lexer) lexStateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(tokNeq)
		return lexText
	}

	l.emit(tokNot)

	return lexText
}

func lexLt(l *lexer) lexStateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(tokLte)
		return lexText
	}

	l.emit(tokLt)

	return lexText
}

func lexGt(l *lexer) lexStateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(tokGte)
		return lexText
	}

	l.emit(tokGt)

	return lexText
}

func lexNumber(l *lexer) lexStateFn {
	l.backup()

	l.accept("+-")
	l.acceptRun("0123456789")
	if l.accept(".") {
		l.acceptRun("0123456789")
	}

	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}

	l.emit(tokNumber)

	return lexText
}

func lexString(l *lexer) lexStateFn {
loop:
	for {
		switch l.next() {
		case '"':
			break loop
		}
	}

	l.emit(tokString)

	return lexText
}

func lexIdentifier(l *lexer) lexStateFn {
loop:
	for {
		switch ch := l.next(); {
		case isAlphaNumeric(ch):
			// absorb
		default:
			l.backup()
			// TODO(vincent): error checking at terminator
			word := l.input[l.start:l.pos]
			switch {
			case word == "where":
				l.emit(tokWhere)
			case word == "and":
				l.emit(tokAnd)
			case word == "or":
				l.emit(tokOr)
			case word == "in":
				l.emit(tokIn)
			case word == "contains":
				l.emit(tokContains)
			case word == "true", word == "false":
				l.emit(tokBool)
			default:
				l.emit(tokIdentifier)
			}
			break loop
		}
	}
	return lexText
}

func lexSpace(l *lexer) lexStateFn {
	for isWhitespace(l.peek()) {
		l.next() // absorb
	}
	l.ignore()
	return lexText
}

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isAlphaNumeric(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch)
}
