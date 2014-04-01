package php

import (
	"testing"

	"stephensearles.com/php/token"
)

var testFile = `<?php

$str = "I'm a string!";
$num = -230.098;

session_start();

function foo(barType $bar, $foobar) {
  if (buzz()) {
    fizzbuzz();
  }
}

class MyClass {
  public myMethod() {
  }
}

$myvar = $this->myMethod();

?>
<html>
<? echo something(); ?>
</html>`

func assertNext(t *testing.T, l *lexer, typ token.Token) Item {
	i := l.nextItem()
	if i.typ != typ {
		t.Fatal("Incorrect lexing. Expected:", typ, "Found:", i)
	}
	return i
}

func assertItem(t *testing.T, i Item, expected string) {
	if i.val != expected {
		t.Fatal("Did not correctly parse token.", i)
	}
}

func TestPHPLexer(t *testing.T) {
	l := newLexer(testFile)

	var i Item
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

	i = assertNext(t, l, token.Class)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.BlockBegin)
	i = assertNext(t, l, token.Public)
	i = assertNext(t, l, token.Identifier)
	i = assertNext(t, l, token.OpenParen)
	i = assertNext(t, l, token.CloseParen)
	i = assertNext(t, l, token.BlockBegin)
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
