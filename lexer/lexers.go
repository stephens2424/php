package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/stephens2424/php/token"
)

// longestToken is length of the longest token string
var longestToken = 0

const shortPHPBegin = "<?"
const longPHPBegin = "<?php"
const phpEnd = "?>"

const eof = -1

func init() {
	for k := range token.TokenMap {
		if len(k) > longestToken {
			longestToken = len(k)
		}
	}
}

// lexHTML consumes and emits an html t until it
// finds a php begin
func lexHTML(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], shortPHPBegin) {
			if l.pos > l.start {
				l.emit(token.HTML)
			}
			return lexPHPBegin
		}
		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(token.HTML)
	}
	l.emit(token.EOF)
	return nil
}

func lexPHPBegin(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], longPHPBegin) {
		l.pos += len(longPHPBegin)
	}
	if strings.HasPrefix(l.input[l.pos:], shortPHPBegin) {
		l.pos += len(shortPHPBegin)
	}
	l.emit(token.PHPBegin)
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
		l.emit(token.EOF)
		return nil
	}

	// strings
	if l.peek() == '`' {
		return lexShellCommand
	}

	if l.peek() == '\'' {
		return lexSingleQuotedStringLiteral
	}

	if l.peek() == '"' {
		return lexDoubleQuotedStringLiteral
	}

	if t, ok := hasKeyword(l); ok {
		l.emit(t)
		return lexPHP
	}

	return lexIdentifier
}

func isOperator(t token.Token) bool {
	_, ok := map[token.Token]struct{}{
		token.VariableOperator: struct{}{},
		token.ObjectOperator: struct{}{},
		token.ScopeResolutionOperator: struct{}{},
	}[t]

	return ok
}

func hasKeyword(l *lexer) (token.Token, bool) {
	var t token.Token
	// is it an operator or keyword?
	//
	// start by getting the longest possible keyword string
	// out of the remaining input. find the longest keyword
	// match in the input
	var tokenString string
	if len(l.input[l.pos:]) > longestToken {
		tokenString = l.input[l.pos : l.pos+longestToken]
	} else {
		tokenString = l.input[l.pos:]
	}

	tokenString = strings.ToLower(tokenString)

	// starting with the longest match, iterate to see if we have found
	// a keyword or operator token
	for ; tokenString != ""; tokenString = tokenString[:len(tokenString)-1] {
		t, ok := token.TokenMap[tokenString]
		if !ok {
			continue
		}

		if !IsKeyword(t, tokenString) {
			l.pos += len(tokenString)
			return t, true
		}

		if (isOperator(l.getPrevious().Typ)) {
			// if the keyword is preceded by a variable
			// operator, object operator, or scope resolution
			// operator, we actually have an identifier.
			return t, false
		}

		// we think we're at a token of some kind
		l.pos += len(tokenString)
		if l.accept(alphabet + underscore + digits) {
			// but if the keyword actually continues on
			// unexpectedly, roll back because this is
			// actually an identifier

			// to account for the extra character consumed by
			// accept
			l.backup()

			// move back the length of the false keyword now
			l.pos -= len(tokenString)
			return t, false
		}

		// we definitely have a token, emit it.
		return t, true
	}
	return t, false
}

func lexIdentifier(l *lexer) stateFn {
	var weirdIdentifier bool
	l.acceptRunFn(func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '\\' {
			return true
		}

		if unicode.IsSpace(r) {
			return false
		}

		if _, ok := token.OperatorMarks[r]; ok {
			return false
		}

		weirdIdentifier = true
		return true
	})

	if weirdIdentifier {
		l.errorf("unexpected characters in identifier")
	}

	l.emit(token.Identifier)
	return lexPHP
}

func lexNumberLiteral(l *lexer) stateFn {
	if l.accept("0") {
		// binary?
		if l.accept("b") {
			l.acceptRun("01")
			l.emit(token.NumberLiteral)
			return lexPHP
		}
		// hexadecimal?
		if l.accept("xX") {
			l.acceptRun(digits + "abcdefABCDEF")
			l.emit(token.NumberLiteral)
			return lexPHP
		}
	}
	// is decimal?
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.accept("eE") {
		l.acceptRun(digits)
	}

	l.emit(token.NumberLiteral)
	return lexPHP
}

func lexShellCommand(l *lexer) stateFn {
	l.next()
	for {
		switch l.next() {
		case '`':
			l.emit(token.ShellCommand)
			return lexPHP
		case eof:
			l.emit(token.ShellCommand)
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
			l.emit(token.StringLiteral)
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
			l.emit(token.StringLiteral)
			return lexPHP
		}
	}
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"
const underscore = "_"

// lexPHPEnd lexes the end of a PHP section returning the context to HTML
func lexPHPEnd(l *lexer) stateFn {
	l.pos += len(phpEnd)
	l.emit(token.PHPEnd)
	return lexHTML
}

func lexLineComment(l *lexer) stateFn {
	lineLength := strings.IndexAny(l.input[l.pos:], "\r\n") + 1
	if lineLength == 0 {
		// this is the last line, so lex until the end
		lineLength = len(l.input[l.pos:])
	}

	// deal with varying line endings
	if l.input[l.pos+lineLength-1] == '\r' &&
		len(l.input[l.pos+lineLength:]) > 0 &&
		l.input[l.pos+lineLength] == '\n' {
		lineLength++
	}

	// don't lex php end
	if phpEndLength := strings.Index(l.input[l.pos:l.pos+lineLength], phpEnd); phpEndLength >= 0 && phpEndLength < lineLength {
		lineLength = phpEndLength
	}
	l.pos += lineLength
	l.emit(token.CommentLine)
	return lexPHP
}

func lexBlockComment(l *lexer) stateFn {
	commentLength := strings.Index(l.input[l.pos:], "*/") + 2
	if commentLength == 1 {
		// the file ends before we find */
		commentLength = len(l.input[l.pos:])
	}
	l.pos += commentLength
	l.emit(token.CommentBlock)
	return lexPHP
}

func lexDoc(l *lexer) stateFn {
	var nowDoc bool
	l.pos += len("<<<")
	l.skipSpace()
	if strings.HasPrefix(l.input[l.pos:], "'") {
		nowDoc = true
		l.pos += len("'")
	} else if l.peek() == '"' {
		l.next()
	}
	labelPos := l.pos
	l.accept(underscore + alphabet)
	l.acceptRun(underscore + alphabet + digits)

	endMarkerA := fmt.Sprintf("\r\n%s", l.input[labelPos:l.pos])
	endMarkerB := fmt.Sprintf("\n%s", l.input[labelPos:l.pos])
	endMarkerC := fmt.Sprintf("\r%s", l.input[labelPos:l.pos])
	if nowDoc {
		l.accept("'")
	} else if l.peek() == '"' {
		l.next()
	}
	l.accept("\r\n")
	l.accept("\r\n")

	for !strings.HasPrefix(l.input[l.pos:], endMarkerA) &&
		!strings.HasPrefix(l.input[l.pos:], endMarkerB) &&
		!strings.HasPrefix(l.input[l.pos:], endMarkerC) {
		l.next()
	}

	if strings.HasPrefix(l.input[l.pos:], endMarkerA) {
		l.pos += len(endMarkerA)
	} else {
		l.pos += len(endMarkerB)
	}

	l.emit(token.StringLiteral)
	return lexPHP
}
