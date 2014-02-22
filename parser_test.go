package php

import (
	"fmt"
	"reflect"
	"testing"

	"stephensearles.com/php/ast"
)

func TestPHPParserHW(t *testing.T) {
	testStr := `hello world`
	p := newParser(testStr)
	a := p.parse()
	if len(a) != 1 || a[0] != ast.Echo(ast.Literal{Type: ast.String}) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestPHPParserHWPHP(t *testing.T) {
	testStr := `<?php
    echo "hello world";`
	p := newParser(testStr)
	a := p.parse()
	if len(a) != 1 || !reflect.DeepEqual(a[0], ast.Echo(ast.Literal{Type: ast.String})) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestIf(t *testing.T) {
	testStr := `<?php
    if (true)
      echo "hello world";
    else if (false)
      echo "no hello world";`
	p := newParser(testStr)
	a := p.parse()
	ifStmtOne := ast.IfStmt{
		Condition: ast.Literal{Type: ast.Boolean},
		TrueBlock: ast.Echo(ast.Literal{Type: ast.String}),
		FalseBlock: &ast.IfStmt{
			Condition:  ast.Literal{Type: ast.Boolean},
			TrueBlock:  ast.Echo(ast.Literal{Type: ast.String}),
			FalseBlock: ast.Block{},
		},
	}
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok := a[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !reflect.DeepEqual(*parsedIf, ifStmtOne) {
		t.Fatalf("If did not correctly parse")
	}
}

func TestAssignment(t *testing.T) {
	testStr := `<?php
    $test = "hello world";
    echo $test;`
	p := newParser(testStr)
	a := p.parse()
	if len(a) != 2 {
		t.Fatalf("Assignment did not correctly parse")
	}
}

func TestFunction(t *testing.T) {
	testStr := `<?php
    function TestFn($arg) {
      echo $arg;
    }
    $var = TestFn("world");`
	p := newParser(testStr)
	a := p.parse()
	if len(a) != 2 {
		t.Fatalf("Function did not correctly parse")
	}
	_, ok := a[0].(*ast.FunctionStmt)
	if !ok {
		t.Fatalf("FunctionStmt did not correctly parse")
	}
	assign, ok := a[1].(ast.AssignmentStmt)
	if !ok {
		t.Fatalf("FunctionCall did not correctly parse")
	}
	_, ok = assign.Value.(*ast.FunctionCallExpression)
	if !ok {
		t.Fatalf("FunctionCall did not correctly parse")
	}
}

func TestClass(t *testing.T) {
	testStr := `<?php
    class TestClass {
      public function method1($arg) {
        echo $arg;
      }
      private function method2(TestClass $arg) {
        echo $arg;
        return $arg;
      }
    }`
	p := newParser(testStr)
	a := p.parse()
	if len(a) != 1 {
		t.Fatalf("Class did not correctly parse")
	}
	parsedClass, ok := a[0].(ast.Class)
	if !ok {
		t.Fatalf("Class did not correctly parse")
	}
	if parsedClass.Name != "TestClass" {
		t.Fatalf("Class Name did not correctly parse. Got:%s", parsedClass.Name)
	}
	if len(parsedClass.Methods) != 2 {
		t.Fatalf("Class methods did not correctly parse")
	}
	if parsedClass.Methods[0].Name != "method1" {
		t.Fatalf("Class method did not correctly parse. Got:%s", parsedClass.Methods[0].Name)
	}
	if parsedClass.Methods[1].Name != "method2" {
		t.Fatalf("Class method did not correctly parse. Got:%s", parsedClass.Methods[0].Name)
	}
	if parsedClass.Methods[1].Arguments[0].TypeHint != "TestClass" {
		t.Fatalf("Class method did not correctly parse. Got:%s", parsedClass.Methods[0].Name)
	}
}

func TestExpressionParsing(t *testing.T) {
	p := newParser(`<? if (1 + 2 > 3)
    echo "good"; `)
	a := p.parse()
	ifStmt := ast.IfStmt{
		Condition: ast.OperatorExpression{
			Operand1: ast.OperatorExpression{
				Operand1: ast.Literal{Type: ast.Float},
				Operand2: ast.Literal{Type: ast.Float},
				Type:     ast.Numeric,
				Operator: "+",
			},
			Operand2: ast.Literal{Type: ast.Float},
			Type:     ast.Boolean,
			Operator: ">",
		},
		TrueBlock:  ast.Echo(ast.Literal{Type: ast.String}),
		FalseBlock: ast.Block{},
	}
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok := a[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !reflect.DeepEqual(*parsedIf, ifStmt) {
		t.Fatalf("If did not correctly parse")
	}

	p = newParser(`<? if (4 + 5 * 6)
    echo "bad";
  `)
	a = p.parse()
	ifStmt = ast.IfStmt{
		Condition: ast.OperatorExpression{
			Operand2: ast.OperatorExpression{
				Operand1: ast.Literal{Type: ast.Float},
				Operand2: ast.Literal{Type: ast.Float},
				Type:     ast.Numeric,
				Operator: "*",
			},
			Operand1: ast.Literal{Type: ast.Float},
			Type:     ast.Numeric,
			Operator: "+",
		},
		TrueBlock:  ast.Echo(ast.Literal{Type: ast.String}),
		FalseBlock: ast.Block{},
	}
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok = a[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !reflect.DeepEqual(*parsedIf, ifStmt) {
		t.Fatalf("If did not correctly parse")
	}

	p = newParser(`<? if (1 > 2 * 3 + 4)
    echo "good";
  `)
	a = p.parse()
	ifStmt = ast.IfStmt{
		Condition: ast.OperatorExpression{
			Operand1: ast.Literal{Type: ast.Float},
			Operand2: ast.OperatorExpression{
				Operand1: ast.OperatorExpression{
					Operand1: ast.Literal{Type: ast.Float},
					Operand2: ast.Literal{Type: ast.Float},
					Type:     ast.Numeric,
					Operator: "*",
				},
				Operand2: ast.Literal{Type: ast.Float},
				Operator: "+",
				Type:     ast.Numeric,
			},
			Type:     ast.Boolean,
			Operator: ">",
		},
		TrueBlock:  ast.Echo(ast.Literal{Type: ast.String}),
		FalseBlock: ast.Block{},
	}
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok = a[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !reflect.DeepEqual(*parsedIf, ifStmt) {
		t.Fatalf("If did not correctly parse")
	}

	p = newParser(`<? if (1 > 2 * 3 + 4 - 2 & 3 && 4 ^ 8 or 14 xor 10 and 13 >> 18 << 10)
    echo "good";
  `)
	p.debug = true
	a = p.parse()
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
}

func TestArray(t *testing.T) {
	testStr := `<?
  $var = array("one", "two", "three");`
	p := newParser(testStr)
	p.debug = true
	a := p.parse()
	if len(a) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.AssignmentStmt{
		Assignee: ast.NewIdentifier("$var"),
		Operator: "=",
		Value: &ast.ArrayExpression{
			ast.BaseNode{},
			ast.ArrayType{},
			[]ast.ArrayPair{
				{nil, ast.Literal{Type: ast.String}},
				{nil, ast.Literal{Type: ast.String}},
				{nil, ast.Literal{Type: ast.String}},
			},
		},
	}
	if !reflect.DeepEqual(a[0], tree) {
		fmt.Printf("Found:    %+v\n", a[0])
		fmt.Printf("Expected: %+v\n", tree)
		t.Fatalf("Array did not correctly parse")
	}
}

func TestArrayKeys(t *testing.T) {
	testStr := `<?
  $var = array(1 => "one", 2 => "two", 3 => "three");`
	p := newParser(testStr)
	a := p.parse()
	if len(a) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.AssignmentStmt{
		Assignee: ast.NewIdentifier("$var"),
		Operator: "=",
		Value: &ast.ArrayExpression{
			ast.BaseNode{},
			ast.ArrayType{},
			[]ast.ArrayPair{
				{ast.Literal{Type: ast.Float}, ast.Literal{Type: ast.String}},
				{ast.Literal{Type: ast.Float}, ast.Literal{Type: ast.String}},
				{ast.Literal{Type: ast.Float}, ast.Literal{Type: ast.String}},
			},
		},
	}
	if !reflect.DeepEqual(a[0], tree) {
		fmt.Printf("Found:    %+v\n", a[0])
		fmt.Printf("Expected: %+v\n", tree)
		t.Fatalf("Array did not correctly parse")
	}
}

func TestMethodCall(t *testing.T) {
	testStr := `<?
  $res = $var->do();`
	p := newParser(testStr)
	p.debug = true
	a := p.parse()
	if len(a) == 0 {
		t.Fatalf("Method call did not correctly parse")
	}
}
