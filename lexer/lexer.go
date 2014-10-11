// Lexer transforms an input string into a stream of PHP tokens.
package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/stephens2424/php/token"
)

// Lexer represents the state of "lexing" items from a source string.
// The idea is derived from a Rob Pike talk:
// http://www.youtube.com/watch?v=HxaD_trXwRE
type lexer struct {
	// start stores the start position of the currently lexing token..
	start int
	// lastStart stores the start position of the previously lexed token..
	lastStart int
	// lastPos stores the position of the previous lexed element.
	lastPos int

	// pos is the current position of the lexer in the input, as an index
	// of the input string.
	pos  int
	line int
	// width is the length of the current rune
	width int
	items chan token.Item // channel of scanned items.

	// input is the full input string.
	input string

	// file is the filename of the input, used to print errors.
	file string
}

func NewLexer(input string) token.Stream {
	l := &lexer{
		line:  1,
		input: input,
		items: make(chan token.Item),
	}
	go l.run()
	return l
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// Run lexes the input by executing state functions until
// the state is nil. It is typically called in a goroutine.
func (l *lexer) run() {
	for state := lexHTML; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// emit gets the current token., sends it on the token. channel
// and prepares for lexing the next token.
func (l *lexer) emit(t token.Token) {
	i := token.Item{
		Typ:   t,
		Begin: l.currentLocation(),
		Val:   l.input[l.start:l.pos],
	}

	l.incrementLines()
	l.lastPos = i.Position().Position
	l.start = l.pos

	i.End = l.currentLocation()
	l.items <- i
}

func (l *lexer) currentLocation() token.Position {
	return token.Position{Position: l.start, Line: l.line, File: l.file}
}

// nextItem returns the next token. from the input.
func (l *lexer) Next() token.Item {
	Item := <-l.items
	return Item
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) previous() rune {
	r, _ := utf8.DecodeRuneInString(l.input[l.lastPos:])
	return r
}

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

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
	l.incrementLines()
}

func (l *lexer) skipSpace() {
	r := l.next()
	for isSpace(r) {
		r = l.next()
	}
	l.backup()
	l.ignore()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	i := token.Item{
		Typ:   token.Error,
		Begin: l.currentLocation(),
		End:   l.currentLocation(),
		Val:   fmt.Sprintf(format, args...),
	}
	l.incrementLines()
	l.items <- i
	return nil
}

func (l *lexer) incrementLines() {
	l.line += strings.Count(l.input[l.lastStart:l.pos], "\n")
	l.lastStart = l.pos
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func IsKeyword(i token.Token, tokenString string) bool {
	_, ok := keywordMap[i]
	return ok && !isNonAlphaOperator(tokenString)
}

var nonalpha *regexp.Regexp

func init() {
	nonalpha = regexp.MustCompile(`^[^a-zA-Z0-9]*$`)
}

func isNonAlphaOperator(s string) bool {
	return nonalpha.MatchString(s)
}

// keywordMap lists all keywords that should be ignored as a prefix to a longer
// identifier.
var keywordMap = map[token.Token]bool{}

func init() {
	re := regexp.MustCompile("^[a-zA-Z]+")
	for keyword, t := range token.TokenMap {
		if re.MatchString(keyword) {
			keywordMap[t] = true
		}
	}
}
