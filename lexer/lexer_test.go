package lexer

import (
	"testing"

	"github.com/stephens2424/php/token"
)

var testFile = `<?php

$str = "I'm a string!";
$num = -230.098;

session_start();

/** comment */
function foo(barType $bar, $foobar) {
  if (buzz()) {
    fizzbuzz();
  }
}

// comment
class MyClass {
  public myMethod() {
    ++self::$num;
  }
}

$myvar = $this->myMethod();

$x =<<<"foo"
    bar
foo;

?>
<html>
<? echo something(); ?>
</html>`

func assertNext(t *testing.T, l token.Stream, typ token.Token) token.Item {
	i := l.Next()
	if i.Typ != typ {
		t.Fatal("Incorrect lexing. Expected:", typ, "Found:", i)
	}
	return i
}

func assertItem(t *testing.T, i token.Item, expected string) {
	if i.Val != expected {
		t.Fatal("Did not correctly parse token.", i)
	}
}

func TestPHPLexer(t *testing.T) {
	l := token.Subset(NewLexer(testFile), token.Significant|token.CommentType)

	var i token.Item
	i = assertNext(t, l, token.PHPBegin)

	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.AssignmentOperator)
	i = assertNext(t, l, token.StringLiteral)
	i = assertNext(t, l, token.StatementEnd)

	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.AssignmentOperator)
	i = assertNext(t, l, token.SubtractionOperator)
	i = assertNext(t, l, token.NumberLiteral)
	i = assertNext(t, l, token.StatementEnd)

	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.StatementEnd)

	i = assertNext(t, l, token.CommentBlock)
	i = assertNext(t, l, token.Function)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.Comma)
	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.BlockBegin)

	i = assertNext(t, l, token.If)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.CloseParen)

	i = assertNext(t, l, token.BlockBegin)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.StatementEnd)
	i = assertNext(t, l, token.BlockEnd)

	i = assertNext(t, l, token.BlockEnd)

	i = assertNext(t, l, token.CommentLine)
	i = assertNext(t, l, token.Class)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.BlockBegin)
	i = assertNext(t, l, token.Public)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.BlockBegin)
	i = assertNext(t, l, token.UnaryOperator)
	i = assertNext(t, l, token.Self)
	i = assertNext(t, l, token.ScopeResolutionOperator)
	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.StatementEnd)
	i = assertNext(t, l, token.BlockEnd)
	i = assertNext(t, l, token.BlockEnd)

	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.AssignmentOperator)
	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.ObjectOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.StatementEnd)

	i = assertNext(t, l, token.VariableOperator)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.AssignmentOperator)
	i = assertNext(t, l, token.StringLiteral)
	i = assertNext(t, l, token.StatementEnd)

	i = assertNext(t, l, token.PHPEnd)
	i = assertNext(t, l, token.HTML)
	assertItem(t, i, "\n<html>\n")

	i = assertNext(t, l, token.PHPBegin)
	i = assertNext(t, l, token.Echo)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.StatementEnd)

	i = assertNext(t, l, token.PHPEnd)
	i = assertNext(t, l, token.HTML)
	assertItem(t, i, "\n</html>")

	i = assertNext(t, l, token.EOF)
}
