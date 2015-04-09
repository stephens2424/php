package deadcode

import (
	"testing"

	"github.com/stephens2424/php"
	"github.com/stephens2424/php/ast"
)

func TestDeadCode(t *testing.T) {
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

	p := php.NewParser()
	_, errs := p.Parse("test.php", src)
	if len(errs) > 0 {
		t.Fatal(errs)
	}

	var shouldBeDead = map[string]struct{}{
		"simple": struct{}{},
		"other":  struct{}{},
	}

	dead := DeadFunctions(p.FileSet, []string{"test.php"})

	for _, deadFunc := range dead {

		fnName := deadFunc.(*ast.FunctionStmt).Name
		if _, ok := shouldBeDead[fnName]; !ok {
			t.Error("%q was found dead, but shouldn't have been", fnName)
		}
		delete(shouldBeDead, fnName)
	}

	for fugitive, _ := range shouldBeDead {
		t.Error("%q should have been found dead, but wasn't", fugitive)
	}
}
