package php

import (
	"fmt"
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

	if strings.HasPrefix(l.input[l.pos:], "<<<") {
		return lexDoc
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

	if l.peek() == eof {
		l.emit(itemEOF)
		return nil
	}

	if l.peek() == '`' {
		return lexShellCommand
	}

	if l.peek() == '\'' {
		return lexSingleQuotedStringLiteral
	}

	if l.peek() == '"' {
		return lexDoubleQuotedStringLiteral
	}

	for _, token := range tokenList {
		item := tokenMap[token]
		potentialToken := l.input[l.pos:]
		if len(potentialToken) > len(token) {
			potentialToken = potentialToken[:len(token)]
		}
		if strings.HasPrefix(strings.ToLower(potentialToken), token) {
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

	l.acceptRun(alphabet + underscore + digits + "\\")
	l.emit(itemIdentifier)
	return lexPHP
}

func lexNumberLiteral(l *lexer) stateFn {
	if l.accept("0") {
		if l.accept("x") {
			// hexadecimal
			l.acceptRun(digits + "abcdef")
			l.emit(itemNumberLiteral)
			return lexPHP
		}
	}
	// is decimal?
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(itemNumberLiteral)
	return lexPHP
}

func lexShellCommand(l *lexer) stateFn {
	l.next()
	for {
		switch l.next() {
		case '`':
			l.emit(itemShellCommand)
			return lexPHP
		case eof:
			l.emit(itemShellCommand)
			return nil
		}
	}
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

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"
const underscore = "_"

func lexIdentifier(l *lexer) stateFn {
	l.accept("$")
	l.accept(underscore + alphabet)
	l.acceptRun(underscore + alphabet + digits)
	l.emit(itemVariableOperator)
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

func lexDoc(l *lexer) stateFn {
	var nowDoc bool
	l.pos += len("<<<")
	if strings.HasPrefix(l.input[l.pos:], "'") {
		nowDoc = true
	}
	labelPos := l.pos
	l.accept(underscore + alphabet)
	l.acceptRun(underscore + alphabet + digits)
	endMarker := fmt.Sprintf("\n%s", l.input[labelPos:l.pos])
	if nowDoc {
		l.accept("'")
	}
	l.accept("\n")
	for !strings.HasPrefix(l.input[l.pos:], endMarker) {
		l.next()
	}
	l.pos += len(endMarker)
	l.emit(itemStringLiteral)
	return lexPHP
}
