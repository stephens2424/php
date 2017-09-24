// Package lexer transforms an input string into a stream of PHP tokens.
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
	start     int // start stores the start position of the currently lexing token..
	lastStart int // lastStart stores the start position of the previously lexed token..
	lastPos   int // lastPos stores the position of the previous lexed element.

	abort   chan struct{}
	pos     int             // pos is the current position of the lexer in the input, as an index of the input string.
	line    int             // line is the current line number
	width   int             // width is the length of the current rune
	itemsCh chan token.Item // channel of scanned items.
	items   []token.Item    // the items lexed so far
	itemPos int             // the current position in items

	// input is the full input string.
	input string

	// file is the filename of the input, used to print errors.
	file string

	lastSignificant token.Item
}

// NewLexer returns a token stream
func NewLexer(input string) token.Stream {
	l := &lexer{
		line:    1,
		input:   input,
		itemsCh: make(chan token.Item),
		abort:   make(chan struct{}),
	}
	go l.run()
	return l
}

func (l *lexer) Abort() {
	close(l.abort)
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// Run lexes the input by executing state functions until
// the state is nil. It is typically called in a goroutine.
func (l *lexer) run() {
	for state := lexHTML; state != nil; {
		state = state(l)
		select {
		case <-l.abort:
			goto done
		default:
		}
	}
done:
	close(l.itemsCh) // No more tokens will be delivered.
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

	select {
	case l.itemsCh <- i:
	case <-l.abort:
	}

	if i.Typ.Type().Is(token.Significant) {
		l.lastSignificant = i
	}
}

func (l *lexer) currentLocation() token.Position {
	return token.Position{Position: l.start, Line: l.line, File: l.file}
}

// Next returns the next token from the input and advances the lexer by a token.
func (l *lexer) Next() token.Item {
	// if we've lexed at least one item and the most recent lexed item is EOF, return the zero value
	if l.itemPos > 0 && l.items[l.itemPos-1].Typ == token.EOF {
		return token.Item{}
	}

	// if Previous has been called and we have already-lexed items pending, return the next of those
	if l.itemPos < len(l.items)-1 {
		item := l.items[l.itemPos]
		l.itemPos++
		return item
	}

	// lex a new item and return it
	item := <-l.itemsCh
	l.items = append(l.items, item)
	l.itemPos++
	return item
}

// getPrevious returns the most recently lexed significant item without
// modifying the lexer state.
func (l *lexer) getPrevious() token.Item {
	return l.lastSignificant
}

// Previous returns the most recently lexed item and resets the position
// of the lexer back by an item.
func (l *lexer) Previous() token.Item {
	// if we have no previous items, return the zero value
	if l.itemPos <= 0 {
		return token.Item{}
	}

	l.itemPos--
	return l.items[l.itemPos]
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

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// acceptRunFn consumes a run of runes from the valid set function.
func (l *lexer) acceptRunFn(validFn func(r rune) bool) {
	for validFn(l.next()) {
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

func (l *lexer) skipSpace() {
	r := l.next()
	for isSpace(r) {
		l.emit(token.Space)
		r = l.next()
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	i := token.Item{
		Typ:   token.Error,
		Begin: l.currentLocation(),
		End:   l.currentLocation(),
		Val:   fmt.Sprintf(format, args...),
	}
	l.itemsCh <- i
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

// IsKeyword returns whether token is a keyword
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
