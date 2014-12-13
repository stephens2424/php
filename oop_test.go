package php

import (
	"testing"

	"github.com/stephens2424/php/ast"
)

func TestClass(t *testing.T) {
	testStr := `<?php
    abstract class TestClass {
      public $myProp;
      protected $myProp2;
      const my_const = "test";
      private $arr = array("one", "two");
      abstract public function method0($arg);
      public function method1($arg) {
        echo $arg;
      }
      private function method2(TestClass $arg, $arg2 = false) {
        echo $arg;
        return $arg;
      }
    }`
	p := NewParser(testStr)
	p.Debug = true
	a, errs := p.Parse()
	if len(errs) > 0 {
		t.Fatal(errs)
	}
	if len(a) != 1 {
		t.Fatalf("Class did not correctly parse")
	}
	tree := &ast.Class{
		Name: "TestClass",
		Constants: []ast.Constant{
			{
				Variable: ast.NewVariable("my_const"),
				Value:    &ast.Literal{Type: ast.String, Value: `"test"`},
			},
		},
		Properties: []ast.Property{
			{
				Visibility: ast.Public,
				Name:       "$myProp",
			},
			{
				Visibility: ast.Protected,
				Name:       "$myProp2",
			},
			{
				Visibility: ast.Private,
				Name:       "$arr",
				Initialization: &ast.ArrayExpression{
					Pairs: []ast.ArrayPair{
						{Value: &ast.Literal{Type: ast.String, Value: `"one"`}},
						{Value: &ast.Literal{Type: ast.String, Value: `"two"`}},
					},
				},
			},
		},
		Methods: []ast.Method{
			{
				Visibility: ast.Public,
				FunctionStmt: &ast.FunctionStmt{
					FunctionDefinition: &ast.FunctionDefinition{
						Name: "method0",
						Arguments: []ast.FunctionArgument{
							{
								Variable: ast.NewVariable("arg"),
							},
						},
					},
				},
			},
			{
				Visibility: ast.Public,
				FunctionStmt: &ast.FunctionStmt{
					FunctionDefinition: &ast.FunctionDefinition{
						Name: "method1",
						Arguments: []ast.FunctionArgument{
							{
								Variable: ast.NewVariable("arg"),
							},
						},
					},
					Body: &ast.Block{
						Statements: []ast.Statement{
							ast.Echo(ast.NewVariable("arg")),
						},
					},
				},
			},
			{
				Visibility: ast.Private,
				FunctionStmt: &ast.FunctionStmt{
					FunctionDefinition: &ast.FunctionDefinition{
						Name: "method2",
						Arguments: []ast.FunctionArgument{
							{
								TypeHint: "TestClass",
								Variable: ast.NewVariable("arg"),
							},
							{
								Variable: ast.NewVariable("arg2"),
								Default:  &ast.Literal{Type: ast.Boolean, Value: "false"},
							},
						},
					},
					Body: &ast.Block{
						Statements: []ast.Statement{
							ast.Echo(ast.NewVariable("arg")),
							ast.ReturnStmt{Expression: ast.NewVariable("arg")},
						},
					},
				},
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Class did not parse correctly")
	}
}

func TestExtraModifiers(t *testing.T) {
	testStr := `<?
  class myclass {
    public public function test() {
    }
  }`

	p := NewParser(testStr)
	_, errs := p.Parse()
	if len(errs) != 1 {
		t.Fatalf("Did not correctly error that a function has two public modifiers")
	}
}

func TestInstantiation(t *testing.T) {
	testStr := `<?
  $obj = new Obj::$classes['obj']($arg);`
	p := NewParser(testStr)
	a, errs := p.Parse()
	if len(errs) != 0 {
		t.Fatalf("Did not parse instantiation correctly: %s", errs)
	}
	tree := ast.ExpressionStmt{ast.AssignmentExpression{
		Operator: "=",
		Assignee: ast.NewVariable("obj"),
		Value: &ast.NewExpression{
			Class: ast.NewClassExpression("Obj", &ast.ArrayLookupExpression{
				Array: ast.NewVariable("classes"),
				Index: &ast.Literal{Type: ast.String, Value: `'obj'`},
			}),
			Arguments: []ast.Expression{
				ast.NewVariable("arg"),
			},
		},
	}}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Instantiation did not parse correctly")
	}
}
