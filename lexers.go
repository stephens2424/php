package php

import (
	"strings"
	"unicode"
)

const shortPHPBegin = "<?"
const longPHPBegin = "<?php"
const phpEnd = "?>"

const eof = -1

// lexHTML consumes and emits an html item until it
// finds a php begin
func lexHTML(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], shortPHPBegin) {
			if l.pos > l.start {
				l.emit(itemHTML)
			}
			return lexPHPBegin
		}
		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(itemHTML)
	}
	l.emit(itemEOF)
	return nil
}

func lexPHPBegin(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], longPHPBegin) {
		l.pos += len(longPHPBegin)
	}
	if strings.HasPrefix(l.input[l.pos:], shortPHPBegin) {
		l.pos += len(shortPHPBegin)
	}
	l.emit(itemPHPBegin)
	return lexPHP
}

func lexPHP(l *lexer) stateFn {
	l.skipSpace()

	if r := l.peek(); unicode.IsDigit(r) {
		return lexNumberLiteral
	} else if r == '.' {
		l.next() // must advance because we only peeked before
		secondR := l.peek()
		l.backup()
		if unicode.IsDigit(secondR) {
			return lexNumberLiteral
		}
	}

	if strings.HasPrefix(l.input[l.pos:], "?>") {
		return lexPHPEnd
	}

	if strings.HasPrefix(l.input[l.pos:], "#") {
		return lexLineComment
	}

	if strings.HasPrefix(l.input[l.pos:], "//") {
		return lexLineComment
	}

	if strings.HasPrefix(l.input[l.pos:], "/*") {
		return lexBlockComment
	}

	for _, token := range tokenList {
		item := tokenMap[token]
		if strings.HasPrefix(l.input[l.pos:], token) {
			l.pos += len(token)
			if isKeyword(item) && l.accept(alphabet+underscore+digits) {
				l.backup() // to account for the character consumed by accept
				l.pos -= len(token)
				break
			}
			l.emit(item)
			return lexPHP
		}
	}

	if strings.HasPrefix(l.input[l.pos:], "$") {
		return lexIdentifier
	}

	if l.peek() == eof {
		l.emit(itemEOF)
		return nil
	}

	if l.peek() == '\'' {
		return lexSingleQuotedStringLiteral
	}

	if l.peek() == '"' {
		return lexDoubleQuotedStringLiteral
	}

	l.acceptRun(alphabet + underscore + digits + "\\")
	l.emit(itemNonVariableIdentifier)
	return lexPHP
}

func lexNumberLiteral(l *lexer) stateFn {
	// is decimal?
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(itemNumberLiteral)
	return lexPHP
}

func lexSingleQuotedStringLiteral(l *lexer) stateFn {
	l.next()
	for {
		switch l.next() {
		case '\\':
			l.next()
			continue
		case '\'':
			l.emit(itemStringLiteral)
			return lexPHP
		}
	}
}

func lexDoubleQuotedStringLiteral(l *lexer) stateFn {
	l.next()
	for {
		switch l.next() {
		case '\\':
			l.next()
			continue
		case '"':
			l.emit(itemStringLiteral)
			return lexPHP
		}
	}
}

func lexIf(l *lexer) stateFn {
	return l.errorf("if is not supported")
}

func lexCondition(l *lexer) stateFn {
	// this could be useful in the condition of a while, do-while, for terminator, if, and if-else block
	// what state should it return?
	// in all cases except do-while, after this is done, a block-begin is the correct state

	// how can this take advantage of the lexPHP function?
	return lexPHP
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"
const underscore = "_"

func lexIdentifier(l *lexer) stateFn {
	l.accept("$")
	l.accept(underscore + alphabet)
	l.acceptRun(underscore + alphabet + digits)
	l.emit(itemIdentifier)
	return lexPHP
}

// lexBlockBegin lexes the beginning of a code block delimited by '{'.
// This state occurs after the declaration of control flow structures.
func lexBlockBegin(l *lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	if l.next() == '{' {
		l.emit(itemBlockBegin)
	} else {
		l.errorf("expecting { to begin a new block")
	}
	return lexPHP
}

// lexBlockEnd lexes the end of a code block delimited by '}'.
func lexBlockEnd(l *lexer) stateFn {
	l.pos += 1
	l.emit(itemBlockEnd)
	return lexPHP
}

// lexPHPEnd lexes the end of a PHP section returning the context to HTML
func lexPHPEnd(l *lexer) stateFn {
	l.pos += len(phpEnd)
	l.emit(itemPHPEnd)
	return lexHTML
}

func lexLineComment(l *lexer) stateFn {
	lineLength := strings.Index(l.input[l.pos:], "\n") + 1
	if lineLength == 0 {
		// this is the last line, so lex until the end
		lineLength = len(l.input[l.pos:])
	}
	// don't lex php end
	if phpEndLength := strings.Index(l.input[l.pos:l.pos+lineLength], phpEnd); phpEndLength >= 0 && phpEndLength < lineLength {
		lineLength = phpEndLength
	}
	l.pos += lineLength
	l.ignore()
	return lexPHP
}

func lexBlockComment(l *lexer) stateFn {
	commentLength := strings.Index(l.input[l.pos:], "*/") + 2
	if commentLength == 1 {
		// the file ends before we find */
		commentLength = len(l.input[l.pos:])
	}
	if phpEndLength := strings.Index(l.input[l.pos:l.pos+commentLength], phpEnd); phpEndLength >= 0 {
		commentLength = phpEndLength
	}
	l.pos += commentLength
	l.ignore()
	return lexPHP
}
