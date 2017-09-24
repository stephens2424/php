package deadcode

import (
	"testing"

	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/parser"
)

func TestDeadFunctions(t *testing.T) {
	src := `<?php
	$var1 = "a";
	function simple() {
		$var2 = "b";
		$var3 = "c";
	}

	class fizz {
		const buzz = "fizzbuzz";

		static function notsimple() {
			$var4 = "d";
		}

		function other() {}
	}

	fizz::notsimple();
	`

	p := parser.NewParser()
	if _, err := p.Parse("test.php", src); err != nil {
		t.Fatal(err)
	}

	var shouldBeDead = map[string]struct{}{
		"simple": struct{}{},
		"other":  struct{}{},
	}

	dead := DeadFunctions(p.FileSet, []string{"test.php"})

	for _, deadFunc := range dead {

		fnName := deadFunc.(*ast.FunctionStmt).Name
		if _, ok := shouldBeDead[fnName]; !ok {
			t.Errorf("%q was found dead, but shouldn't have been", fnName)
		}
		delete(shouldBeDead, fnName)
	}

	for _, fugitive := range shouldBeDead {
		t.Errorf("%q should have been found dead, but wasn't", fugitive)
	}
}
