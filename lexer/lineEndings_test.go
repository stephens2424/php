package lexer

import (
	"fmt"
	"testing"

	"github.com/stephens2424/php/token"
)

func TestLineEndings(t *testing.T) {
	var testFiles = []string{
		"<?php\r//comment\r$test;\r",
		"<?php\n//comment\n$test;\n",
		"<?php\r\n//comment\r\n$test;\r\n",
		"<?php\n//comment\r$test;\r",
		"<?php\n//comment\r$test;\n",
	}

	for i, testFile := range testFiles {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			l := token.Subset(NewLexer(testFile), token.Significant|token.CommentType)
			assertNext(t, l, token.PHPBegin)
			assertNext(t, l, token.CommentLine)
			assertNext(t, l, token.VariableOperator)
			assertNext(t, l, token.Identifier)
			assertNext(t, l, token.StatementEnd)
			assertNext(t, l, token.EOF)
		})
	}
}
