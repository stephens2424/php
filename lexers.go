package php

import (
	"strings"
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
	if strings.HasPrefix(l.input[l.pos:], "$") {
		return lexIdentifier
	}
	if strings.HasPrefix(l.input[l.pos:], "function") {
		return lexFunctionDeclaration
	}
	if strings.HasPrefix(l.input[l.pos:], "class") {
		return l.errorf("classes are unsupported")
	}
	if strings.HasPrefix(l.input[l.pos:], "interface") {
		return l.errorf("interfaces are unsupported")
	}
	if strings.HasPrefix(l.input[l.pos:], "{") {
		return lexBlockBegin
	}
	if strings.HasPrefix(l.input[l.pos:], "}") {
		return lexBlockEnd
	}
	if strings.HasPrefix(l.input[l.pos:], ";") {
		l.next()
		l.emit(itemStatementEnd)
		return lexPHP
	}
	if strings.HasPrefix(l.input[l.pos:], "?>") {
		return lexPHPEnd
	}
	if strings.HasPrefix(l.input[l.pos:], "echo") {
		l.pos += len("echo")
		l.emit(itemEcho)
		return lexPHP
	}
	l.accept(alphabet + underscore)
	l.acceptRun(alphabet + underscore + digits)
	if l.peek() == '(' {
		l.emit(itemFunctionName)
		return lexFunctionArgs
	}
	for {
		if l.next() == eof {
			break
		}
	}
	if l.start > l.pos {
		l.emit(itemPHP)
	}
	l.emit(itemEOF)
	return nil
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

func lexFunctionDeclaration(l *lexer) stateFn {
	l.pos += len("function")
	l.emit(itemFunction)
	return lexFunctionName
}

func lexFunctionName(l *lexer) stateFn {
	l.skipSpace()
	l.accept(underscore + alphabet)
	l.acceptRun(underscore + alphabet + digits)
	l.emit(itemFunctionName)
	return lexFunctionArgs
}

func lexFunctionArgs(l *lexer) stateFn {
	l.skipSpace()
	switch r := l.next(); {
	case r == '(':
		l.emit(itemArgumentListBegin)
		if l.peek() == ')' {
			return lexFunctionArgs
		}
		return lexFunctionArg
	case r == ')':
		l.emit(itemArgumentListEnd)
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

// lexOperator lexes any operator
func lexOperator(l *lexer) stateFn {
	l.acceptRun("!*()%<>-=+/")
	l.emit(itemOperator)
	return lexPHP
}

// lexPHPEnd lexes the end of a PHP section returning the context to HTML
func lexPHPEnd(l *lexer) stateFn {
	l.pos += len(phpEnd)
	l.emit(itemPHPEnd)
	return lexHTML
}
