package deadcode

import (
	"testing"

	"github.com/stephens2424/php"
	"github.com/stephens2424/php/ast"
)

func TestDeadClass(t *testing.T) {
	src := `<?php

	class fizz {
		static function a() {}
	}

	class buzz {
		static function b() {}
	}

	class fizzbuzz {
	}

	fizz::notsimple();
	$x = new fizzbuzz();
	`

	p := php.NewParser()
	if _, err := p.Parse("test.php", src); err != nil {
		t.Fatal(err)
	}

	var shouldBeDead = map[string]struct{}{
		"buzz": struct{}{},
	}

	dead := DeadClasses(p.FileSet, []string{"test.php"})

	for _, deadFunc := range dead {
		fnName := deadFunc.(*ast.Class).Name
		if _, ok := shouldBeDead[fnName]; !ok {
			t.Errorf("%q was found dead, but shouldn't have been", fnName)
		}
		delete(shouldBeDead, fnName)
	}

	for fugitive, _ := range shouldBeDead {
		t.Errorf("%q should have been found dead, but wasn't", fugitive)
	}
}
