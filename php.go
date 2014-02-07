package php

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	lastPos int
	pos     int
	start   int
	width   int
	input   string
	items   chan item // channel of scanned items.
}

func newLexer(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}
	return l
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
	for state := lexHTML; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

func (l *lexer) emit(t itemType) {
	i := item{t, l.start, l.input[l.start:l.pos]}
	fmt.Println(i)
	l.items <- i
	l.start = l.pos
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
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
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}
