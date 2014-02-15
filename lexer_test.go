package php

import "testing"

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

?>
<html>
<? echo something(); ?>
</html>`

func assertNext(t *testing.T, l *lexer, typ itemType) item {
	i := l.nextItem()
	if i.typ != typ {
		t.Fatal("Incorrect lexing. Expected:", typ, "Found:", i)
	}
	return i
}

func assertItem(t *testing.T, i item, expected string) {
	if i.val != expected {
		t.Fatal("Did not correctly parse item", i)
	}
}

func TestPHPLexer(t *testing.T) {
	l := newLexer(testFile)

	var i item
	i = assertNext(t, l, itemPHPBegin)

	i = assertNext(t, l, itemIdentifier)
	i = assertNext(t, l, itemAssignmentOperator)
	i = assertNext(t, l, itemStringLiteral)
	i = assertNext(t, l, itemStatementEnd)

	i = assertNext(t, l, itemIdentifier)
	i = assertNext(t, l, itemAssignmentOperator)
	i = assertNext(t, l, itemNumberLiteral)
	i = assertNext(t, l, itemStatementEnd)

	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemOpenParen)
	i = assertNext(t, l, itemCloseParen)
	i = assertNext(t, l, itemStatementEnd)

	i = assertNext(t, l, itemFunction)
	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemOpenParen)
	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemIdentifier)
	i = assertNext(t, l, itemArgumentSeparator)
	i = assertNext(t, l, itemIdentifier)
	i = assertNext(t, l, itemCloseParen)
	i = assertNext(t, l, itemBlockBegin)

	i = assertNext(t, l, itemIf)
	i = assertNext(t, l, itemOpenParen)
	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemOpenParen)
	i = assertNext(t, l, itemCloseParen)
	i = assertNext(t, l, itemCloseParen)

	i = assertNext(t, l, itemBlockBegin)
	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemOpenParen)
	i = assertNext(t, l, itemCloseParen)
	i = assertNext(t, l, itemStatementEnd)
	i = assertNext(t, l, itemBlockEnd)

	i = assertNext(t, l, itemBlockEnd)

	i = assertNext(t, l, itemClass)
	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemBlockBegin)
	i = assertNext(t, l, itemPublic)
	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemOpenParen)
	i = assertNext(t, l, itemCloseParen)
	i = assertNext(t, l, itemBlockBegin)
	i = assertNext(t, l, itemBlockEnd)
	i = assertNext(t, l, itemBlockEnd)

	i = assertNext(t, l, itemPHPEnd)
	i = assertNext(t, l, itemHTML)
	assertItem(t, i, "\n<html>\n")

	i = assertNext(t, l, itemPHPBegin)
	i = assertNext(t, l, itemEcho)
	i = assertNext(t, l, itemNonVariableIdentifier)
	i = assertNext(t, l, itemOpenParen)
	i = assertNext(t, l, itemCloseParen)
	i = assertNext(t, l, itemStatementEnd)

	i = assertNext(t, l, itemPHPEnd)
	i = assertNext(t, l, itemHTML)
	assertItem(t, i, "\n</html>")

	i = assertNext(t, l, itemEOF)
}
