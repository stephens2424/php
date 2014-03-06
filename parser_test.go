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
	a, _ := p.Parse()
	tree := ast.Echo(ast.Literal{Type: ast.String})
	if !assertEquals(a[0], tree) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestPHPParserHWPHP(t *testing.T) {
	testStr := `<?php
    echo "hello world", "!";`
	p := NewParser(testStr)
	a, _ := p.Parse()
	tree := ast.Echo(
		&ast.Literal{Type: ast.String},
		&ast.Literal{Type: ast.String},
	)
	if !assertEquals(a[0], tree) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestInclude(t *testing.T) {
	testStr := `<?php
  include test.php; ?>`
	p := NewParser(testStr)
	_, errs := p.Parse()
	if len(errs) > 0 {
		t.Fatalf("Did not parse include correctly")
	}
}

func TestIf(t *testing.T) {
	testStr := `<?php
    if (true)
      echo "hello world";
    else if (false)
      echo "no hello world";`
	p := NewParser(testStr)
	a, _ := p.Parse()
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
	a, _ := p.Parse()
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
	a, _ := p.Parse()
	tree := []ast.Node{
		&ast.FunctionStmt{
			FunctionDefinition: &ast.FunctionDefinition{
				Name: "TestFn",
				Arguments: []ast.FunctionArgument{
					{
						Variable: ast.NewVariable("$arg"),
					},
				},
			},
			Body: &ast.Block{
				Statements: []ast.Statement{ast.Echo(ast.NewVariable("$arg"))},
			},
		},
		ast.AssignmentStmt{
			ast.AssignmentExpression{
				Assignee: &ast.Variable{Name: "$var", Type: ast.AnyType},
				Value: &ast.FunctionCallExpression{
					FunctionName: ast.FunctionNameExpression{Name: "TestFn"},
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

func TestExpressionParsing(t *testing.T) {
	p := NewParser(`<? if (1 + 2 > 3)
    echo "good"; `)
	a, _ := p.Parse()
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
	a, _ = p.Parse()
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
	a, _ = p.Parse()
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

	p = NewParser(`<? if ($var = &$var2 > 2 * (3 + 4) - 2 & 3 && 4 ^ 8 or 14 xor 10 and 13 >> 18 << 10 ? true : false)
    echo "good";
  `)
	p.Debug = true
	a, _ = p.Parse()
	if len(a) != 1 {
		t.Fatalf("If did not correctly parse")
	}
}

func TestArray(t *testing.T) {
	testStr := `<?
  $var = array("one", "two", "three");`
	p := NewParser(testStr)
	p.Debug = true
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.AssignmentStmt{
		ast.AssignmentExpression{
			Assignee: ast.NewVariable("$var"),
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: ast.NewVariable("$var"),
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Method call did not correctly parse")
	}
	tree := ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: ast.NewVariable("$res"),
		Operator: "=",
		Value: &ast.MethodCallExpression{
			Receiver: ast.NewVariable("$var"),
			FunctionCallExpression: &ast.FunctionCallExpression{
				FunctionName: ast.FunctionNameExpression{Name: "go"},
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
	a, _ := p.Parse()
	if len(a) != 2 {
		t.Fatalf("Property did not correctly parse")
	}
	tree := ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: ast.NewVariable("$res"),
		Operator: "=",
		Value: &ast.PropertyExpression{
			Receiver: ast.NewVariable("$var"),
			Name:     ast.PropertyIdentifier{Name: "go"},
		},
	}}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Property did not correctly parse")
	}

	tree = ast.AssignmentStmt{ast.AssignmentExpression{
		Assignee: &ast.PropertyExpression{
			Receiver: ast.NewVariable("$var"),
			Name:     ast.PropertyIdentifier{Name: "go"},
		},
		Operator: "=",
		Value:    ast.NewVariable("$res"),
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Do loop did not correctly parse")
	}
	tree := &ast.DoWhileStmt{
		Termination: ast.NewVariable("$otherVar"),
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("$var")),
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.WhileStmt{
		Termination: ast.NewVariable("$otherVar"),
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("$var")),
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.ForeachStmt{
		Source: &ast.Variable{Name: "$arr", Type: ast.AnyType},
		Key:    &ast.Variable{Name: "$key", Type: ast.AnyType},
		Value:  &ast.Variable{Name: "$val", Type: ast.AnyType},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{ast.Echo(ast.OperatorExpression{
				Operator: ".",
				Operand1: &ast.Variable{Name: "$key", Type: ast.AnyType},
				Operand2: &ast.Variable{Name: "$val", Type: ast.AnyType},
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.ForStmt{
		Initialization: ast.AssignmentExpression{
			Assignee: &ast.Variable{Type: ast.AnyType, Name: "$i"},
			Value:    &ast.Literal{Type: ast.Float},
			Operator: "=",
		},
		Termination: ast.OperatorExpression{
			Operand1: ast.NewVariable("$i"),
			Operand2: &ast.Literal{Type: ast.Float},
			Operator: "<",
			Type:     ast.Boolean,
		},
		Iteration: ast.OperatorExpression{
			Operator: "++",
			Operand1: ast.NewVariable("$i"),
			Type:     ast.Numeric,
		},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("$i")),
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.WhileStmt{
		Termination: ast.AssignmentExpression{
			Assignee: ast.NewVariable("$var"),
			Value: &ast.FunctionCallExpression{
				FunctionName: ast.FunctionNameExpression{Name: "mysql_assoc"},
				Arguments:    make([]ast.Expression, 0),
			},
			Operator: "=",
		},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("$var")),
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("While loop with assignment did not correctly parse")
	}
}

func TestArrayLookup(t *testing.T) {
	testStr := `<?
  echo $arr['one'][$two];
  $var->arr[] = 2;
  echo $arr[2 + 1];`
	p := NewParser(testStr)
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array lookup did not correctly parse")
	}
	tree := []ast.Node{
		ast.EchoStmt{
			Expressions: []ast.Expression{&ast.ArrayLookupExpression{
				Array: &ast.ArrayLookupExpression{
					Array: &ast.Variable{Name: "$arr", Type: ast.AnyType},
					Index: &ast.Literal{Type: ast.String},
				},
				Index: &ast.Variable{Name: "$two", Type: ast.AnyType},
			}},
		},
		ast.AssignmentStmt{
			ast.AssignmentExpression{
				Assignee: ast.ArrayAppendExpression{
					Array: &ast.PropertyExpression{
						Receiver: ast.NewVariable("$var"),
						Name:     ast.PropertyIdentifier{Name: "arr"},
					},
				},
				Operator: "=",
				Value:    &ast.Literal{Type: ast.Float},
			},
		},
	}
	if !assertEquals(a[0], tree[0]) {
		t.Fatalf("Array lookup did not correctly parse")
	}
	if !assertEquals(a[1], tree[1]) {
		t.Fatalf("Array append expression did not correctly parse")
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
	a, _ := p.Parse()
	if len(a) == 0 {
		t.Fatalf("Array lookup did not correctly parse")
	}
	tree := ast.SwitchStmt{
		Expression: &ast.Variable{Name: "$var", Type: ast.AnyType},
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
	a, _ := p.Parse()
	if len(a) != 4 {
		t.Fatalf("Literals did not correctly parse")
	}
	tree := []ast.Node{
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewVariable("$var"),
			Value:    &ast.Literal{Type: ast.String},
			Operator: "=",
		}},
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewVariable("$var"),
			Value:    &ast.Literal{Type: ast.Float},
			Operator: "=",
		}},
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewVariable("$var"),
			Value:    &ast.Literal{Type: ast.Boolean},
			Operator: "=",
		}},
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: ast.NewVariable("$var"),
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
		ast.Echo(ast.Literal{Type: ast.String}),
	}
	p := NewParser(testStr)
	a, _ := p.Parse()
	if !reflect.DeepEqual(a, tree) {
		fmt.Printf("Found:    %+v\n", a)
		fmt.Printf("Expected: %+v\n", tree)
		t.Fatalf("Literals did not correctly parse")
	}
}

func TestScopeResolutionOperator(t *testing.T) {
	testStr := `<?
  MyClass::myfunc($var);
  echo MyClass::myconst;
  echo $var::myfunc();`
	p := NewParser(testStr)
	a, _ := p.Parse()
	tree := []ast.Node{
		ast.ExpressionStmt{
			&ast.ClassExpression{
				Receiver: "MyClass",
				Expression: &ast.FunctionCallExpression{
					FunctionName: ast.FunctionNameExpression{Name: "myfunc"},
					Arguments: []ast.Expression{
						&ast.Variable{Name: "$var", Type: ast.AnyType},
					},
				},
			},
		},
		ast.Echo(&ast.ClassExpression{
			Receiver: "MyClass",
			Expression: ast.ConstantExpression{
				&ast.Variable{Name: "myconst", Type: ast.AnyType},
			},
		}),
		ast.Echo(&ast.ClassExpression{
			Receiver: "$var",
			Expression: &ast.FunctionCallExpression{
				FunctionName: ast.FunctionNameExpression{Name: "myfunc"},
				Arguments:    []ast.Expression{},
			},
		}),
	}
	if !assertEquals(a[0], tree[0]) {
		t.Fatalf("Scope resolution operator function call did not correctly parse")
	}
	if !assertEquals(a[1], tree[1]) {
		t.Fatalf("Scope resolution operator expression did not correctly parse")
	}
	if !assertEquals(a[2], tree[2]) {
		t.Fatalf("Scope resolution operator function call on identifier did not correctly parse")
	}
}

func TestCastOperator(t *testing.T) {
	testStr := `<?
  $var = (double) 1.0; ?>`
	p := NewParser(testStr)
	a, _ := p.Parse()
	tree := []ast.Node{
		ast.AssignmentStmt{ast.AssignmentExpression{
			Assignee: &ast.Variable{Name: "$var", Type: ast.AnyType},
			Value: ast.OperatorExpression{
				Operand1: &ast.Literal{Type: ast.Float},
				Operator: "(double)",
				Type:     ast.Numeric,
			},
			Operator: "=",
		}},
	}
	if !assertEquals(a[0], tree[0]) {
		t.Fatalf("Cast operator parsing failed")
	}
}

func TestInterface(t *testing.T) {
	testStr := `<?
  interface MyInterface extends YourInterface, HerInterface {
    public function TheirFunc();
    private function MyFunc();
  }`
	p := NewParser(testStr)
	a, _ := p.Parse()
	tree := &ast.Interface{
		Name:     "MyInterface",
		Inherits: []string{"YourInterface", "HerInterface"},
		Methods: []ast.Method{
			{
				Visibility: ast.Public,
				FunctionStmt: &ast.FunctionStmt{
					FunctionDefinition: &ast.FunctionDefinition{
						Name:      "TheirFunc",
						Arguments: []ast.FunctionArgument{},
					},
				},
			},
			{
				Visibility: ast.Private,
				FunctionStmt: &ast.FunctionStmt{
					FunctionDefinition: &ast.FunctionDefinition{
						Name:      "MyFunc",
						Arguments: []ast.FunctionArgument{},
					},
				},
			},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Interface did not parse correctly")
	}
}

func TestGlobal(t *testing.T) {
	testStr := `<?
  global $var, $otherVar;`
	p := NewParser(testStr)
	a, _ := p.Parse()
	tree := &ast.GlobalDeclaration{
		Identifiers: []*ast.Variable{
			{Name: "$var", Type: ast.AnyType},
			{Name: "$otherVar", Type: ast.AnyType},
		},
	}
	if !assertEquals(a[0], tree) {
		t.Fatalf("Global did not parse correctly")
	}
}
