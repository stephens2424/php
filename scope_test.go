package php

import (
	"testing"

	"github.com/stephens2424/php/ast"
)

func TestScope(t *testing.T) {
	src := `<?php
	$var1 = "a";
	function simple() {
		$var2 = "b";
		$var3 = "c";
	}

	class fizz {
		const buzz = "fizzbuzz";

		function notsimple() {
			$var4 = "d";
		}
	}
	`

	p := NewParser()
	a, errs := p.Parse("test.php", src)
	if len(errs) > 0 {
		t.Fatal(errs)
	}

	ExpectFunctions(a.Namespace, []string{"simple"}, t)
	ExpectClasses(a.Namespace, []string{"fizz"}, t)
}

func ExpectFunctions(ns *ast.Namespace, functions []string, t *testing.T) {
	expected := map[string]struct{}{}
	hasError := false
	for _, fn := range functions {
		expected[fn] = struct{}{}
		if _, ok := ns.Functions[fn]; !ok {
			t.Errorf("expected function %q, but didn't find it", fn)
			hasError = true
		}
	}

	for fn := range ns.Functions {
		if _, ok := expected[fn]; !ok {
			t.Errorf("found function %q, but didn't expect it", fn)
			hasError = true
		}
	}

	if hasError {
		t.FailNow()
	}
}

func ExpectClasses(ns *ast.Namespace, classes []string, t *testing.T) {
	expected := map[string]struct{}{}
	hasError := false
	for _, class := range classes {
		expected[class] = struct{}{}
		if _, ok := ns.ClassesAndInterfaces[class]; !ok {
			t.Errorf("expected class %q, but didn't find it", class)
			hasError = true
		}
	}

	for class := range ns.ClassesAndInterfaces {
		if _, ok := expected[class]; !ok {
			t.Errorf("found class %q, but didn't expect it", class)
			hasError = true
		}
	}

	if hasError {
		t.FailNow()
	}
}
