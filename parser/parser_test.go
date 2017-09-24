package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/passes/printing"
)

func assertEquals(found, expected ast.Node) bool {
	w := printing.NewWalker()
	if !reflect.DeepEqual(found, expected) {
		fmt.Printf("Found:    %s\n", found)
		w.Walk(found)
		fmt.Printf("Expected: %+s\n", expected)
		w.Walk(expected)
		findDifference(found, expected)
		return false
	}
	return true
}

func findDifference(found, expected ast.Node) {
	w := printing.NewWalker()
	foundChildren := found.Children()
	expectedChildren := expected.Children()
	if len(foundChildren) != len(expectedChildren) {
		fmt.Printf("Found Subtree:    %s\n", found)
		w.Walk(found)
		fmt.Printf("Expected Subtree: %+s\n", expected)
		w.Walk(expected)
	} else if len(foundChildren) != 0 && len(foundChildren) == len(expectedChildren) {
		for i := 0; i < len(foundChildren); i++ {
			if !reflect.DeepEqual(foundChildren[i], expectedChildren[i]) {
				findDifference(foundChildren[i], expectedChildren[i])
			}
		}
	} else {
		fmt.Printf("Found Subtree:    %s\n", found)
		w.Walk(found)
		fmt.Printf("Expected Subtree: %+s\n", expected)
		w.Walk(expected)
	}
}

func TestPHPParserHW(t *testing.T) {
	testStr := `hello world`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := ast.Echo(ast.Literal{Type: ast.String, Value: `hello world`})
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestNestedIfElse(t *testing.T) {
	testStr := `
<?php function foo (){if (true) {} else {} someExpr(); }
`
	p := NewParser()
	p.disableScoping = true
	ast, err := p.Parse("", testStr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if ast.Nodes[1].Children()[1].Children()[1].String() != "someExpr()" {
		t.Fatal("ast incorrect")
	}
}

func TestPHPParserHWPHP(t *testing.T) {
	testStr := `<?php
    echo "hello world", "!";`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := ast.Echo(
		&ast.Literal{Type: ast.String, Value: `"hello world"`},
		&ast.Literal{Type: ast.String, Value: `"!"`},
	)
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Hello world did not correctly parse")
	}
}

func TestInclude(t *testing.T) {
	testStr := `<?php
  include "test.php"; ?>`
	p := NewParser()
	p.disableScoping = true
	_, err := p.Parse("test.php", testStr)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("Did not parse include correctly")
	}
}

func TestIf(t *testing.T) {
	testStr := `<?php
    if (true)
      echo "hello world";
    else if (false)
      echo "no hello world";`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := &ast.IfStmt{
		Branches: []ast.IfBranch{
			{
				Condition: &ast.Literal{Type: ast.Boolean, Value: "true"},
				Block:     ast.Echo(&ast.Literal{Type: ast.String, Value: `"hello world"`}),
			},
			{
				Condition: &ast.Literal{Type: ast.Boolean, Value: "false"},
				Block:     ast.Echo(&ast.Literal{Type: ast.String, Value: `"no hello world"`}),
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("If did not correctly parse")
	}
}

func TestIfBraces(t *testing.T) {
	testStr := `<?php
    if (true) {
      echo "hello world";
    } else if (false) {
      echo "no hello world";
    }`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := &ast.IfStmt{
		Branches: []ast.IfBranch{
			{
				Condition: &ast.Literal{Type: ast.Boolean, Value: "true"},
				Block: &ast.Block{
					Statements: []ast.Statement{ast.Echo(&ast.Literal{Type: ast.String, Value: `"hello world"`})},
				},
			},
			{
				Condition: &ast.Literal{Type: ast.Boolean, Value: "false"},
				Block: &ast.Block{
					Statements: []ast.Statement{ast.Echo(&ast.Literal{Type: ast.String, Value: `"no hello world"`})},
				},
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("If with braces did not correctly parse")
	}
}

func TestAssignment(t *testing.T) {
	testStr := `<?php
    $test = "hello world";
    echo $test;`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) != 2 {
		t.Fatalf("Assignment did not correctly parse")
	}
}

func TestFunctionReturnTypeDeclaration(t *testing.T) {
	testStr := `<?php function TestFn(MyClass $myClass): string {}`
	p := NewParser()
	tree, err := p.Parse("", testStr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	tp := tree.Nodes[0].Children()[0].(*ast.FunctionDefinition).Type
	if tp != "string" {
		t.Fatalf("invalid return type: %v", t)
	}
}

func TestFunction(t *testing.T) {
	testStr := `<?php
    function TestFn($arg) {
      echo $arg;
    }
    $var = TestFn("world", 0);`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := []ast.Node{
		&ast.FunctionStmt{
			FunctionDefinition: &ast.FunctionDefinition{
				Name: "TestFn",
				Arguments: []*ast.FunctionArgument{
					{
						Variable: ast.NewVariable("arg"),
					},
				},
			},
			Body: &ast.Block{
				Statements: []ast.Statement{ast.Echo(ast.NewVariable("arg"))},
			},
		},
		ast.ExprStmt{
			Expr: ast.AssignmentExpr{
				Assignee: ast.NewVariable("var"),
				Value: &ast.FunctionCallExpr{
					FunctionName: &ast.Identifier{Value: "TestFn"},
					Arguments: []ast.Expr{
						&ast.Literal{Type: ast.String, Value: `"world"`},
						&ast.Literal{Type: ast.Float, Value: "0"},
					},
				},
				Operator: "=",
			},
		},
	}
	if len(a.Nodes) != 2 {
		t.Fatalf("Function did not correctly parse")
	}
	if !assertEquals(a.Nodes[0], tree[0]) {
		t.Fatalf("Function did not correctly parse")
	}
	if !assertEquals(a.Nodes[1], tree[1]) {
		t.Fatalf("Function assignment did not correctly parse")
	}
}

func TestExpressionParsing(t *testing.T) {
	p := NewParser()
	p.disableScoping = true
	testStr := `<? if (1 + 2 > 3)
    echo "good"; `
	a, _ := p.Parse("test.php", testStr)
	ifStmt := ast.IfStmt{
		Branches: []ast.IfBranch{
			{
				Condition: ast.BinaryExpr{
					Antecedent: ast.BinaryExpr{
						Antecedent: &ast.Literal{Type: ast.Float, Value: "1"},
						Subsequent: &ast.Literal{Type: ast.Float, Value: "2"},
						Type:       ast.Numeric,
						Operator:   "+",
					},
					Subsequent: &ast.Literal{Type: ast.Float, Value: "3"},
					Type:       ast.Boolean,
					Operator:   ">",
				},
				Block: ast.Echo(&ast.Literal{Type: ast.String, Value: `"good"`}),
			},
		},
	}
	if len(a.Nodes) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok := a.Nodes[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !assertEquals(*parsedIf, ifStmt) {
		t.Fatalf("If did not correctly parse")
	}

	p = NewParser()
	p.disableScoping = true
	testStr = `<? if (4 + 5 * 6)
    echo "bad";
  `
	a, _ = p.Parse("test.php", testStr)
	ifStmt = ast.IfStmt{
		Branches: []ast.IfBranch{
			{
				Condition: ast.BinaryExpr{
					Subsequent: ast.BinaryExpr{
						Antecedent: &ast.Literal{Type: ast.Float, Value: "5"},
						Subsequent: &ast.Literal{Type: ast.Float, Value: "6"},
						Type:       ast.Numeric,
						Operator:   "*",
					},
					Antecedent: &ast.Literal{Type: ast.Float, Value: "4"},
					Type:       ast.Numeric,
					Operator:   "+",
				},
				Block: ast.Echo(&ast.Literal{Type: ast.String, Value: `"bad"`}),
			},
		},
	}
	if len(a.Nodes) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok = a.Nodes[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !reflect.DeepEqual(*parsedIf, ifStmt) {
		t.Fatalf("If did not correctly parse")
	}

	p = NewParser()
	p.disableScoping = true
	testStr = `<? if (1 > 2 * 3 + 4)
    echo "good";
  `
	a, _ = p.Parse("test.php", testStr)
	ifStmt = ast.IfStmt{
		Branches: []ast.IfBranch{
			{
				Condition: ast.BinaryExpr{
					Antecedent: &ast.Literal{Type: ast.Float, Value: `1`},
					Subsequent: ast.BinaryExpr{
						Antecedent: ast.BinaryExpr{
							Antecedent: &ast.Literal{Type: ast.Float, Value: `2`},
							Subsequent: &ast.Literal{Type: ast.Float, Value: `3`},
							Type:       ast.Numeric,
							Operator:   "*",
						},
						Subsequent: &ast.Literal{Type: ast.Float, Value: `4`},
						Operator:   "+",
						Type:       ast.Numeric,
					},
					Type:     ast.Boolean,
					Operator: ">",
				},
				Block: ast.Echo(&ast.Literal{Type: ast.String, Value: `"good"`}),
			},
		},
	}
	if len(a.Nodes) != 1 {
		t.Fatalf("If did not correctly parse")
	}
	parsedIf, ok = a.Nodes[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("If did not correctly parse")
	}
	if !reflect.DeepEqual(*parsedIf, ifStmt) {
		t.Fatalf("If did not correctly parse")
	}

	p = NewParser()
	p.disableScoping = true
	testStr = `<? if ($var = &$var2 > 2 * (3 + 4) - 2 & 3 && 4 ^ 8 or 14 xor 10 and 13 >> 18 << 10 ? true : false)
    echo "good";
  `
	p.Debug = true
	a, _ = p.Parse("test.php", testStr)
	if len(a.Nodes) != 1 {
		t.Fatalf("Expression did not correctly parse")
	}
}

func TestArray(t *testing.T) {
	testStr := `<?
  $var = array("one", "two", "three");`
	p := NewParser()
	p.disableScoping = true
	p.Debug = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.ExprStmt{
		Expr: ast.AssignmentExpr{
			Assignee: ast.NewVariable("var"),
			Operator: "=",
			Value: &ast.ArrayExpr{
				ArrayType: ast.ArrayType{},
				Pairs: []ast.ArrayPair{
					{Value: &ast.Literal{Type: ast.String, Value: `"one"`}},
					{Value: &ast.Literal{Type: ast.String, Value: `"two"`}},
					{Value: &ast.Literal{Type: ast.String, Value: `"three"`}},
				},
			},
		},
	}
	if !reflect.DeepEqual(a.Nodes[0], tree) {
		fmt.Printf("Found:    %+v\n", a.Nodes[0])
		fmt.Printf("Expected: %+v\n", tree)
		t.Fatalf("Array did not correctly parse")
	}
}

func TestArrayKeys(t *testing.T) {
	testStr := `<?
  $var = array(1 => "one", 2 => "two", 3 => "three");`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("Array did not correctly parse")
	}
	tree := ast.ExprStmt{Expr: ast.AssignmentExpr{
		Assignee: ast.NewVariable("var"),
		Operator: "=",
		Value: &ast.ArrayExpr{
			ArrayType: ast.ArrayType{},
			Pairs: []ast.ArrayPair{
				{Key: &ast.Literal{Type: ast.Float, Value: "1"}, Value: &ast.Literal{Type: ast.String, Value: `"one"`}},
				{Key: &ast.Literal{Type: ast.Float, Value: "2"}, Value: &ast.Literal{Type: ast.String, Value: `"two"`}},
				{Key: &ast.Literal{Type: ast.Float, Value: "3"}, Value: &ast.Literal{Type: ast.String, Value: `"three"`}},
			},
		},
	}}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Array did not correctly parse")
	}
}

func TestMethodCall(t *testing.T) {
	testStr := `<?
  $res = $var->go();`
	p := NewParser()
	p.disableScoping = true
	p.Debug = true
	p.MaxErrors = 0
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("Method call did not correctly parse")
	}
	tree := ast.ExprStmt{Expr: ast.AssignmentExpr{
		Assignee: ast.NewVariable("res"),
		Operator: "=",
		Value: &ast.MethodCallExpr{
			Receiver: ast.NewVariable("var"),
			FunctionCallExpr: &ast.FunctionCallExpr{
				FunctionName: &ast.Identifier{Value: "go"},
				Arguments:    make([]ast.Expr, 0),
			},
		},
	}}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Method call did not correctly parse")
	}
}

func TestProperty(t *testing.T) {
	testStr := `<?
  $res = $var->go;
  $var->go = $res;`
	p := NewParser()
	p.disableScoping = true
	p.Debug = true
	p.MaxErrors = 0
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) != 2 {
		t.Fatalf("Property did not correctly parse")
	}
	tree := ast.ExprStmt{Expr: ast.AssignmentExpr{
		Assignee: ast.NewVariable("res"),
		Operator: "=",
		Value: &ast.PropertyCallExpr{
			Receiver: ast.NewVariable("var"),
			Name:     &ast.Identifier{Value: "go"},
		},
	}}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Property did not correctly parse")
	}

	tree = ast.ExprStmt{Expr: ast.AssignmentExpr{
		Assignee: &ast.PropertyCallExpr{
			Receiver: ast.NewVariable("var"),
			Name:     &ast.Identifier{Value: "go"},
		},
		Operator: "=",
		Value:    ast.NewVariable("res"),
	}}
	if !assertEquals(a.Nodes[1], tree) {
		t.Fatalf("Property did not correctly parse")
	}
}

func TestDoLoop(t *testing.T) {
	testStr := `<?
  do {
    echo $var;
  } while ($otherVar);`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("Do loop did not correctly parse")
	}
	tree := &ast.DoWhileStmt{
		Termination: ast.NewVariable("otherVar"),
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("var")),
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("TestLoop did not correctly parse")
	}
}

func TestWhileLoop(t *testing.T) {
	testStr := `<?
  while ($otherVar) {
    echo $var;
  }`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.WhileStmt{
		Termination: ast.NewVariable("otherVar"),
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("var")),
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("TestLoop did not correctly parse")
	}
}

func TestForeachLoop(t *testing.T) {
	testStr := `<?
  foreach ($arr as $key => $val) {
    echo $key . $val;
  } ?>`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.ForeachStmt{
		Source: ast.NewVariable("arr"),
		Key:    ast.NewVariable("key"),
		Value:  ast.NewVariable("val"),
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{ast.Echo(ast.BinaryExpr{
				Operator:   ".",
				Antecedent: ast.NewVariable("key"),
				Subsequent: ast.NewVariable("val"),
				Type:       ast.String,
			})},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Foreach did not correctly parse")
	}
}

func TestForLoop(t *testing.T) {
	testStr := `<?
  for ($i = 0; $i < 10; $i++) {
    echo $i;
  }`
	p := NewParser()
	p.disableScoping = true
	p.Debug = true
	p.MaxErrors = 0
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("For loop did not correctly parse")
	}
	tree := &ast.ForStmt{
		Initialization: []ast.Expr{ast.AssignmentExpr{
			Assignee: ast.NewVariable("i"),
			Value:    &ast.Literal{Type: ast.Float, Value: "0"},
			Operator: "=",
		}},
		Termination: []ast.Expr{ast.BinaryExpr{
			Antecedent: ast.NewVariable("i"),
			Subsequent: &ast.Literal{Type: ast.Float, Value: "10"},
			Operator:   "<",
			Type:       ast.Boolean,
		}},
		Iteration: []ast.Expr{ast.UnaryCallExpr{
			Operator:  "++",
			Operand:   ast.NewVariable("i"),
			Preceding: false,
		}},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("i")),
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("For did not correctly parse")
	}
}

func TestWhileLoopWithAssignment(t *testing.T) {
	testStr := `<?
  while ($var = mysql_assoc()) {
    echo $var;
  }`
	p := NewParser()
	p.disableScoping = true
	p.Debug = true
	p.MaxErrors = 0
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("While loop did not correctly parse")
	}
	tree := &ast.WhileStmt{
		Termination: ast.AssignmentExpr{
			Assignee: ast.NewVariable("var"),
			Value: &ast.FunctionCallExpr{
				FunctionName: &ast.Identifier{Value: "mysql_assoc"},
				Arguments:    make([]ast.Expr, 0),
			},
			Operator: "=",
		},
		LoopBlock: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(ast.NewVariable("var")),
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("While loop with assignment did not correctly parse")
	}
}

func TestArrayLookup(t *testing.T) {
	testStr := `<?
  echo $arr['one'][$two];
  $var->arr[] = 2;
  echo $arr[2 + 1];`
	p := NewParser()
	p.disableScoping = true
	p.Debug = true
	p.MaxErrors = 0
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("Array lookup did not correctly parse")
	}
	tree := []ast.Node{
		ast.EchoStmt{
			Expressions: []ast.Expr{&ast.ArrayLookupExpr{
				Array: &ast.ArrayLookupExpr{
					Array: ast.NewVariable("arr"),
					Index: &ast.Literal{Type: ast.String, Value: `'one'`},
				},
				Index: ast.NewVariable("two"),
			}},
		},
		ast.ExprStmt{
			Expr: ast.AssignmentExpr{
				Assignee: ast.ArrayAppendExpr{
					Array: &ast.PropertyCallExpr{
						Receiver: ast.NewVariable("var"),
						Name:     &ast.Identifier{Value: "arr"},
					},
				},
				Operator: "=",
				Value:    &ast.Literal{Type: ast.Float, Value: "2"},
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree[0]) {
		t.Fatalf("Array lookup did not correctly parse")
	}
	if !assertEquals(a.Nodes[1], tree[1]) {
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
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) == 0 {
		t.Fatalf("Array lookup did not correctly parse")
	}
	tree := ast.SwitchStmt{
		Expr: ast.NewVariable("var"),
		Cases: []*ast.SwitchCase{
			{
				Expr: &ast.Literal{Type: ast.Float, Value: "1"},
				Block: ast.Block{
					Statements: []ast.Statement{
						ast.Echo(&ast.Literal{Type: ast.String, Value: `"one"`}),
					},
				},
			},
			{
				Expr: &ast.Literal{Type: ast.Float, Value: "2"},
				Block: ast.Block{
					Statements: []ast.Statement{
						ast.Echo(&ast.Literal{Type: ast.String, Value: `"two"`}),
					},
				},
			},
		},
		DefaultCase: &ast.Block{
			Statements: []ast.Statement{
				ast.Echo(&ast.Literal{Type: ast.String, Value: `"def"`}),
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Switch did not correctly parse")
	}
}

func TestLiterals(t *testing.T) {
	testStr := `<?
  $var = "one";
  $var = 2;
  $var = true;
  $var = null;`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if len(a.Nodes) != 4 {
		t.Fatalf("Literals did not correctly parse")
	}
	tree := []ast.Node{
		ast.ExprStmt{Expr: ast.AssignmentExpr{
			Assignee: ast.NewVariable("var"),
			Value:    &ast.Literal{Type: ast.String, Value: `"one"`},
			Operator: "=",
		}},
		ast.ExprStmt{Expr: ast.AssignmentExpr{
			Assignee: ast.NewVariable("var"),
			Value:    &ast.Literal{Type: ast.Float, Value: "2"},
			Operator: "=",
		}},
		ast.ExprStmt{Expr: ast.AssignmentExpr{
			Assignee: ast.NewVariable("var"),
			Value:    &ast.Literal{Type: ast.Boolean, Value: "true"},
			Operator: "=",
		}},
		ast.ExprStmt{Expr: ast.AssignmentExpr{
			Assignee: ast.NewVariable("var"),
			Value:    &ast.Literal{Type: ast.Null, Value: "null"},
			Operator: "=",
		}},
	}
	if !reflect.DeepEqual(a.Nodes, tree) {
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
		ast.Echo(ast.Literal{Type: ast.String, Value: "html"}),
	}
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	if !reflect.DeepEqual(a.Nodes, tree) {
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
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := []ast.Node{
		ast.ExprStmt{
			Expr: &ast.ClassExpr{
				Receiver: &ast.Identifier{Value: "MyClass"},
				Expr: &ast.FunctionCallExpr{
					FunctionName: &ast.Identifier{Value: "myfunc"},
					Arguments: []ast.Expr{
						ast.NewVariable("var"),
					},
				},
			},
		},
		ast.Echo(&ast.ClassExpr{
			Receiver: &ast.Identifier{Value: "MyClass"},
			Expr: ast.ConstantExpr{
				Variable: ast.NewVariable("myconst"),
			},
		}),
		ast.Echo(&ast.ClassExpr{
			Receiver: ast.NewVariable("var"),
			Expr: &ast.FunctionCallExpr{
				FunctionName: &ast.Identifier{Value: "myfunc"},
				Arguments:    []ast.Expr{},
			},
		}),
	}
	if !assertEquals(a.Nodes[0], tree[0]) {
		t.Fatalf("Scope resolution operator function call did not correctly parse")
	}
	if !assertEquals(a.Nodes[1], tree[1]) {
		t.Fatalf("Scope resolution operator expression did not correctly parse")
	}
	if !assertEquals(a.Nodes[2], tree[2]) {
		t.Fatalf("Scope resolution operator function call on identifier did not correctly parse")
	}
}

func TestCastOperator(t *testing.T) {
	testStr := `<?
  $var = (double) 1.0; ?>`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := []ast.Node{
		ast.ExprStmt{Expr: ast.AssignmentExpr{
			Assignee: ast.NewVariable("var"),
			Value: ast.UnaryCallExpr{
				Operand:   &ast.Literal{Type: ast.Float, Value: "1.0"},
				Operator:  "(double)",
				Preceding: false,
			},
			Operator: "=",
		}},
	}
	if !assertEquals(a.Nodes[0], tree[0]) {
		t.Fatalf("Cast operator parsing failed")
	}
}

func TestInterface(t *testing.T) {
	testStr := `<?
  interface MyInterface extends YourInterface, HerInterface {
    public function TheirFunc();
    private function MyFunc();
  }`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := &ast.Interface{
		Name:     "MyInterface",
		Inherits: []string{"YourInterface", "HerInterface"},
		Methods: []ast.Method{
			{
				Visibility: ast.Public,
				FunctionStmt: &ast.FunctionStmt{
					FunctionDefinition: &ast.FunctionDefinition{
						Name:      "TheirFunc",
						Arguments: []*ast.FunctionArgument{},
					},
				},
			},
			{
				Visibility: ast.Private,
				FunctionStmt: &ast.FunctionStmt{
					FunctionDefinition: &ast.FunctionDefinition{
						Name:      "MyFunc",
						Arguments: []*ast.FunctionArgument{},
					},
				},
			},
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Interface did not parse correctly")
	}
}

func TestGlobal(t *testing.T) {
	testStr := `<?
  global $var, $otherVar;`
	p := NewParser()
	p.disableScoping = true
	a, _ := p.Parse("test.php", testStr)
	tree := &ast.GlobalDeclaration{
		Identifiers: []*ast.Variable{
			ast.NewVariable("var"),
			ast.NewVariable("otherVar"),
		},
	}
	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Global did not parse correctly")
	}
}
