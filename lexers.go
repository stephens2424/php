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
			l.emit(itemHTML)
			return lexPHPBegin
		}
		if l.next() == eof {
			break
		}
	}
	if l.start > l.pos {
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
		l.pos += len(longPHPBegin)
	}
	l.emit(itemPHPBegin)
	return lexPHP
}

func lexPHP(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], phpEnd) {
			l.emit(itemPHP)
			return lexPHPEnd
		}
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

func lexPHPEnd(l *lexer) stateFn {
	l.pos += len(phpEnd)
	l.emit(itemPHPEnd)
	return lexHTML
}
