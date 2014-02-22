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
	} else if r == '-' {
		l.next()
		if unicode.IsDigit(l.peek()) {
			l.backup()
			return lexNumberLiteral
		}
		if l.peek() == '>' {
			l.next()
			l.emit(itemObjectOperator)
			return lexPHP
		}
		l.emit(itemSubtractionOperator)
		return lexPHP
	}

	if strings.HasPrefix(l.input[l.pos:], "?>") {
		return lexPHPEnd
	}

	for _, token := range tokenList {
		Item := tokenMap[token]
		if strings.HasPrefix(l.input[l.pos:], token) {
			l.pos += len(token)
			l.emit(Item)
			return lexPHP
		}
	}

	if strings.HasPrefix(l.input[l.pos:], "$") {
		return lexIdentifier
	}

	if l.next() == eof {
		l.emit(itemEOF)
		return nil
	}
	l.backup()

	if l.peek() == '\'' {
		return lexSingleQuotedStringLiteral
	}

	if l.peek() == '"' {
		return lexDoubleQuotedStringLiteral
	}

	l.accept(alphabet + underscore + "\\")
	l.acceptRun(alphabet + underscore + digits + "\\")
	l.emit(itemNonVariableIdentifier)
	return lexPHP
}

func lexNumberLiteral(l *lexer) stateFn {
	// is negative?
	l.accept("-")
	l.acceptRun(digits)

	// is decimal?
	if l.accept(".") {
		l.acceptRun(digits)
	}

	l.emit(itemNumberLiteral)
	return lexPHP
}

func lexSingleQuotedStringLiteral(l *lexer) stateFn {
	l.next()
	isEscaped := false
	for {
		r := l.next()
		if r == '\\' {
			isEscaped = true
			continue
		}
		if !isEscaped && r == '\'' {
			break
		}
	}
	l.emit(itemStringLiteral)
	return lexPHP
}

func lexDoubleQuotedStringLiteral(l *lexer) stateFn {
	l.next()
	isEscaped := false
	for {
		r := l.next()
		if r == '\\' {
			isEscaped = true
			continue
		}
		if !isEscaped && r == '"' {
			break
		}
	}
	l.emit(itemStringLiteral)
	return lexPHP
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

func lexFunctionArgs(l *lexer) stateFn {
	l.skipSpace()
	switch r := l.next(); {
	case r == '(':
		l.emit(itemOpenParen)
		if l.peek() == ')' {
			return lexFunctionArgs
		}
		return lexFunctionArg
	case r == ')':
		l.emit(itemCloseParen)
		return lexPHP
	case r == ',':
		l.emit(itemArgumentSeparator)
		return lexFunctionArg
	default:
		return l.errorf("invalid function argument separator: '%s'", string(r))
	}
}

func lexFunctionArg(l *lexer) stateFn {
	l.skipSpace()
	if l.peek() != '$' {
		l.accept(underscore + alphabet)
		l.acceptRun(underscore + alphabet + digits)
		l.emit(itemTypeHint)
	}
	l.skipSpace()
	l.next()
	l.accept(underscore + alphabet)
	l.acceptRun(underscore + alphabet + digits)
	l.emit(itemArgumentName)
	return lexFunctionArgs
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
