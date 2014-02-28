package php

import (
	"fmt"
	"reflect"
	"testing"

	"stephensearles.com/php/ast"
	. "stephensearles.com/php/passes/printing"
)

func assertEquals(found, expected ast.Node) bool {
	w := &Walker{}
	if !reflect.DeepEqual(found, expected) {
		fmt.Printf("Found:    %s\n", found)
		w.Walk(found)
		fmt.Printf("Expected: %+s\n", expected)
		w.Walk(expected)
		return false
	}
	return true
}

func TestPHPParserHW(t *testing.T) {
	testStr := `hello world`
	p := NewParser(testStr)
	a := p.Parse()
	tree := ast.Echo(ast.Literal{Type: ast.String})
	if !assertEquals(a[0], tree) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestPHPParserHWPHP(t *testing.T) {
	testStr := `<?php
    echo "hello world";`
	p := NewParser(testStr)
	a := p.Parse()
	if !assertEquals(a[0], ast.Echo(&ast.Literal{Type: ast.String})) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestIf(t *testing.T) {
	testStr := `<?php
    if (true)
      echo "hello world";
    else if (false)
      echo "no hello world";`
	p := NewParser(testStr)
	a := p.Parse()
	tree := &ast.IfStmt{
		Condition:  &ast.Literal{Type: ast.Boolean},
		TrueBranch: ast.Echo(&ast.Literal{Type: ast.String}),
		FalseBranch: &ast.IfStmt{
			Condition:   &ast.Literal{Type: ast.Boolean},
			TrueBranch:  ast.Echo(&ast.Literal{Type: ast.String}),
			FalseBranch: ast.Block{},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("If did not correctly parse")
	}
}

func TestAssignment(t *testing.T) {
	testStr := `<?php
    $test = "hello world";
    echo $test;`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) != 2 {
		t.Fatalf("Assignment did not correctly parse")
	}
}

func TestFunction(t *testing.T) {
	testStr := `<?php
    function TestFn($arg) {
      echo $arg;
    }
    $var = TestFn("world", 0);`
	p := NewParser(testStr)
	a := p.Parse()
	tree := []ast.Node{
		&ast.FunctionStmt{
			FunctionDefinition: &ast.FunctionDefinition{
				Name: "TestFn",
				Arguments: []ast.FunctionArgument{
					{
						Identifier: ast.NewIdentifier("$arg"),
					},
				},
			},
			Body: &ast.Block{
				Statements: []ast.Statement{ast.Echo(ast.NewIdentifier("$arg"))},
			},
		},
		ast.AssignmentStmt{
			ast.AssignmentExpression{
				Assignee: &ast.Identifier{Name: "$var", Type: ast.AnyType},
				Value: &ast.FunctionCallExpression{
					FunctionName: "TestFn",
					Arguments: []ast.Expression{
						&ast.Literal{Type: ast.String},
						&ast.Literal{Type: ast.Float},
					},
				},
				Operator: "=",
			},
		},
	}
	if len(a) != 2 {
		t.Fatalf("Function did not correctly parse")
	}
	if !assertEquals(a[0], tree[0]) {
		t.Fatalf("Function did not correctly parse")
	}
	if !assertEquals(a[1], tree[1]) {
		t.Fatalf("Function assignment did not correctly parse")
	}
}

func TestClass(t *testing.T) {
	testStr := `<?php
    abstract class TestClass {
      public $myProp;
      protected $myProp2;
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
	a := p.Parse()
	if len(a) != 1 {
		t.Fatalf("Class did not correctly parse")
	}
	tree := ast.Class{
		Name: "TestClass",
		Properties: []ast.Property{
			{
				Visibility: ast.Public,
				Name:       "$myProp",
			},
			{
				Visibility: ast.Protected,
				Name:       "$myProp2",
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
								Identifier: &ast.Identifier{Name: "$arg", Type: ast.AnyType},
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
								Identifier: &ast.Identifier{Name: "$arg", Type: ast.AnyType},
							},
						},
					},
					Body: &ast.Block{
						Statements: []ast.Statement{
							ast.Echo(&ast.Identifier{Name: "$arg", Type: ast.AnyType}),
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
								TypeHint:   "TestClass",
								Identifier: &ast.Identifier{Name: "$arg", Type: ast.AnyType},
							},
							{
								Identifier: &ast.Identifier{Name: "$arg2", Type: ast.AnyType},
								Default:    &ast.Literal{Type: ast.Boolean},
							},
						},
					},
					Body: &ast.Block{
						Statements: []ast.Statement{
							ast.Echo(&ast.Identifier{Name: "$arg", Type: ast.AnyType}),
							ast.ReturnStmt{Expression: &ast.Identifier{Name: "$arg", Type: ast.AnyType}},
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

func TestExpressionParsing(t *testing.T) {
	p := NewParser(`<? if (1 + 2 > 3)
    echo "good"; `)
	a := p.Parse()
	ifStmt := ast.IfStmt{
		Condition: ast.OperatorExpression{
			Operand1: ast.OperatorExpression{
				Operand1: &ast.Literal{Type: ast.Float},
				Operand2: &ast.Literal{Type: ast.Float},
				Type:     ast.Numeric,
				Operator: "+",
			},
			Operand2: &ast.Literal{Type: ast.Float},
			Type:     ast.Boolean,
			Operator: ">",
		},
		TrueBranch:  ast.Echo(&ast.Literal{Type: ast.String}),
		FalseBranch: ast.Block{},
	}
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok := a[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !assertEquals(*parsedIf, ifStmt) {
		t.Fatalf("If did not correctly parse")
	}

	p = NewParser(`<? if (4 + 5 * 6)
    echo "bad";
  `)
	a = p.Parse()
	ifStmt = ast.IfStmt{
		Condition: ast.OperatorExpression{
			Operand2: ast.OperatorExpression{
				Operand1: &ast.Literal{Type: ast.Float},
				Operand2: &ast.Literal{Type: ast.Float},
				Type:     ast.Numeric,
				Operator: "*",
			},
			Operand1: &ast.Literal{Type: ast.Float},
			Type:     ast.Numeric,
			Operator: "+",
		},
		TrueBranch:  ast.Echo(&ast.Literal{Type: ast.String}),
		FalseBranch: ast.Block{},
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

	p = NewParser(`<? if (1 > 2 * 3 + 4)
    echo "good";
  `)
	a = p.Parse()
	ifStmt = ast.IfStmt{
		Condition: ast.OperatorExpression{
			Operand1: &ast.Literal{Type: ast.Float},
			Operand2: ast.OperatorExpression{
				Operand1: ast.OperatorExpression{
					Operand1: &ast.Literal{Type: ast.Float},
					Operand2: &ast.Literal{Type: ast.Float},
					Type:     ast.Numeric,
					Operator: "*",
				},
				Operand2: &ast.Literal{Type: ast.Float},
				Operator: "+",
				Type:     ast.Numeric,
			},
			Type:     ast.Boolean,
			Operator: ">",
		},
		TrueBranch:  ast.Echo(&ast.Literal{Type: ast.String}),
		FalseBranch: ast.Block{},
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

	p = NewParser(`<? if (1 > 2 * 3 + 4 - 2 & 3 && 4 ^ 8 or 14 xor 10 and 13 >> 18 << 10)
    echo "good";
  `)
	p.Debug = true
	a = p.Parse()
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
}

func TestArray(t *testing.T) {
	testStr := `<?
  $var = array("one", "two", "three");`
	p := NewParser(testStr)
	p.Debug = true
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.AssignmentStmt{
		ast.AssignmentExpression{
			Assignee: ast.NewIdentifier("$var"),
			Operator: "=",
			Value: &ast.ArrayExpression{
				ast.BaseNode{},
				ast.ArrayType{},
				[]ast.ArrayPair{
					{Value: &ast.Literal{Type: ast.String}},
					{Value: &ast.Literal{Type: ast.String}},
					{Value: &ast.Literal{Type: ast.String}},
				},
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
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: ast.NewIdentifier("$var"),
		Operator: "=",
		Value: &ast.ArrayExpression{
			ast.BaseNode{},
			ast.ArrayType{},
			[]ast.ArrayPair{
				{Key: &ast.Literal{Type: ast.Float}, Value: &ast.Literal{Type: ast.String}},
				{Key: &ast.Literal{Type: ast.Float}, Value: &ast.Literal{Type: ast.String}},
				{Key: &ast.Literal{Type: ast.Float}, Value: &ast.Literal{Type: ast.String}},
			},
		},
	}}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Array did not correctly parse")
	}
}

func TestMethodCall(t *testing.T) {
	testStr := `<?
  $res = $var->go();`
	p := NewParser(testStr)
	p.Debug = true
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Method call did not correctly parse")
	}
	tree := ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: ast.NewIdentifier("$res"),
		Operator: "=",
		Value: &ast.MethodCallExpression{
			Receiver: ast.NewIdentifier("$var"),
			FunctionCallExpression: &ast.FunctionCallExpression{
				FunctionName: "go",
				Arguments:    make([]ast.Expression, 0),
			},
		},
	}}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Method call did not correctly parse")
	}
}

func TestProperty(t *testing.T) {
	testStr := `<?
  $res = $var->go;
  $var->go = $res;`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) != 2 {
		t.Fatalf("Property did not correctly parse")
	}
	tree := ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: ast.NewIdentifier("$res"),
		Operator: "=",
		Value: &ast.PropertyExpression{
			Receiver: ast.NewIdentifier("$var"),
			Name:     "go",
		},
	}}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Property did not correctly parse")
	}

	tree = ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: &ast.PropertyExpression{
			Receiver: ast.NewIdentifier("$var"),
			Name:     "go",
		},
		Operator: "=",
		Value:    ast.NewIdentifier("$res"),
	}}
	if !assertEquals(a[1], tree) {
		t.Fatalf("Property did not correctly parse")
	}
}

func TestDoLoop(t *testing.T) {
	testStr := `<?
  do {
    echo $var;
  } while ($otherVar);`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Do loop did not correctly parse")
	}
	tree := &ast.DoWhileStmt{
		Termination: ast.NewIdentifier("$otherVar"),
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewIdentifier("$var")),
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("TestLoop did not correctly parse")
	}
}

func TestWhileLoop(t *testing.T) {
	testStr := `<?
  while ($otherVar) {
    echo $var;
  }`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.WhileStmt{
		Termination: ast.NewIdentifier("$otherVar"),
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewIdentifier("$var")),
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("TestLoop did not correctly parse")
	}
}

func TestForeachLoop(t *testing.T) {
	testStr := `<?
  foreach ($arr as $key => $val) {
    echo $key . $val;
  } ?>`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.ForeachStmt{
		Source: &ast.Identifier{Name: "$arr", Type: ast.AnyType},
		Key:    &ast.Identifier{Name: "$key", Type: ast.AnyType},
		Value:  &ast.Identifier{Name: "$val", Type: ast.AnyType},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{ast.Echo(ast.OperatorExpression{
				Operator: ".",
				Operand1: &ast.Identifier{Name: "$key", Type: ast.AnyType},
				Operand2: &ast.Identifier{Name: "$val", Type: ast.AnyType},
				Type:     ast.String,
			})},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Foreach did not correctly parse")
	}
}

func TestForLoop(t *testing.T) {
	testStr := `<?
  for ($i = 0; $i < 10; $i++) {
    echo $i;
  }`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.ForStmt{
		Initialization: ast.AssignmentExpression{
			Assignee: &ast.Identifier{Type: ast.AnyType, Name: "$i"},
			Value:    &ast.Literal{Type: ast.Float},
			Operator: "=",
		},
		Termination: ast.OperatorExpression{
			Operand1: ast.NewIdentifier("$i"),
			Operand2: &ast.Literal{Type: ast.Float},
			Operator: "<",
			Type:     ast.Boolean,
		},
		Iteration: ast.OperatorExpression{
			Operand1: ast.NewIdentifier("$i"),
			Type:     ast.Numeric,
		},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewIdentifier("$i")),
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("TestLoop did not correctly parse")
	}
}

func TestWhileLoopWithAssignment(t *testing.T) {
	testStr := `<?
  while ($var = mysql_assoc()) {
    echo $var;
  }`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.WhileStmt{
		Termination: ast.AssignmentExpression{
			Assignee: ast.NewIdentifier("$var"),
			Value: &ast.FunctionCallExpression{
				FunctionName: "mysql_assoc",
				Arguments:    make([]ast.Expression, 0),
			},
			Operator: "=",
		},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewIdentifier("$var")),
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("While loop with assignment did not correctly parse")
	}
}

func TestArrayLookup(t *testing.T) {
	testStr := `<?
  echo $arr['one'][$two];`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array lookup did not correctly parse")
	}
	tree := ast.EchoStmt{
		Expression: &ast.ArrayLookupExpression{
			Array: &ast.ArrayLookupExpression{
				Array: &ast.Identifier{Name: "$arr", Type: ast.AnyType},
				Index: &ast.Literal{Type: ast.String},
			},
			Index: &ast.Identifier{Name: "$two", Type: ast.AnyType},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Array lookup did not correctly parse")
	}
}

func TestSwitch(t *testing.T) {
	testStr := `<?
  switch ($var) {
  case 1:
    echo "one";
  case 2: {
    echo "two";
  }
  default:
    echo "def";
  }`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array lookup did not correctly parse")
	}
	tree := ast.SwitchStmt{
		Expression: &ast.Identifier{Name: "$var", Type: ast.AnyType},
		Cases: []*ast.SwitchCase{
			{
				Expression: &ast.Literal{Type: ast.Float},
				Block: ast.Block{
					Statements: []ast.Statement{
						ast.Echo(&ast.Literal{Type: ast.String}),
					},
				},
			},
			{
				Expression: &ast.Literal{Type: ast.Float},
				Block: ast.Block{
					Statements: []ast.Statement{
						ast.Echo(&ast.Literal{Type: ast.String}),
					},
				},
			},
		},
		DefaultCase: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(&ast.Literal{Type: ast.String}),
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Switch did not correctly parse")
	}
}

func TestLiterals(t *testing.T) {
	testStr := `<?
  $var = "one";
  $var = 2;
  $var = true;
  $var = null;`
	p := NewParser(testStr)
	a := p.Parse()
	if len(a) != 4 {
		t.Fatalf("Literals did not correctly parse")
	}
	tree := []ast.Node{
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewIdentifier("$var"),
			Value:    &ast.Literal{Type: ast.String},
			Operator: "=",
		}},
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewIdentifier("$var"),
			Value:    &ast.Literal{Type: ast.Float},
			Operator: "=",
		}},
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewIdentifier("$var"),
			Value:    &ast.Literal{Type: ast.Boolean},
			Operator: "=",
		}},
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewIdentifier("$var"),
			Value:    &ast.Literal{Type: ast.Null},
			Operator: "=",
		}},
	}
	if !reflect.DeepEqual(a, tree) {
		fmt.Printf("Found:    %+v\n", a)
		fmt.Printf("Expected: %+v\n", tree)
		t.Fatalf("Literals did not correctly parse")
	}
}

func TestComments(t *testing.T) {
	testStr := `<?
  // comment line
  /*
  block
  */
  #line ?>html`
	tree := []ast.Node{
		ast.EchoStmt{Expression: ast.Literal{Type: ast.String}},
	}
	p := NewParser(testStr)
	a := p.Parse()
	if !reflect.DeepEqual(a, tree) {
		fmt.Printf("Found:    %+v\n", a)
		fmt.Printf("Expected: %+v\n", tree)
		t.Fatalf("Literals did not correctly parse")
	}
}
