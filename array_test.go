package php

import (
	"testing"

	"github.com/stephens2424/php/ast"
)

func TestList(t *testing.T) {
	testStr := `<?
    list($one, $two) = array(1, 2);`

	p := NewParser()
	a, err := p.Parse("test.php", testStr)
	if err != nil {
		t.Fatalf("Did not parse list correctly: %s", err)
	}

	tree := ast.ExpressionStmt{&ast.ListStatement{
		Operator: "=",
		Assignees: []ast.Assignable{
			ast.NewVariable("one"),
			ast.NewVariable("two"),
		},
		Value: &ast.ArrayExpression{
			Pairs: []ast.ArrayPair{
				{Key: nil, Value: &ast.Literal{Value: "1", Type: ast.Float}},
				{Key: nil, Value: &ast.Literal{Value: "2", Type: ast.Float}},
			},
		},
	}}

	if !assertEquals(a.Nodes[0], tree) {
		t.Fatalf("Array bracked did not parse correctly")
	}
}

func TestArrayBracket(t *testing.T) {
	testStr := `<?
    $arr = ["one", "two"];
    $arr2 = ["one" => 1, "two" => 2];`

	p := NewParser()
	a, err := p.Parse("test.php", testStr)
	if err != nil {
		t.Fatalf("Did not parse array bracket correctly: %s", err)
	}

	tree := []ast.Statement{
		ast.ExpressionStmt{ast.AssignmentExpression{
			Operator: "=",
			Assignee: ast.NewVariable("arr"),
			Value: &ast.ArrayExpression{
				Pairs: []ast.ArrayPair{
					{Key: nil, Value: &ast.Literal{Value: `"one"`, Type: ast.String}},
					{Key: nil, Value: &ast.Literal{Value: `"two"`, Type: ast.String}},
				},
			},
		}},
		ast.ExpressionStmt{ast.AssignmentExpression{
			Operator: "=",
			Assignee: ast.NewVariable("arr2"),
			Value: &ast.ArrayExpression{
				Pairs: []ast.ArrayPair{
					{
						Key:   &ast.Literal{Value: `"one"`, Type: ast.String},
						Value: &ast.Literal{Value: "1", Type: ast.Float},
					},
					{
						Key:   &ast.Literal{Value: `"two"`, Type: ast.String},
						Value: &ast.Literal{Value: "2", Type: ast.Float},
					},
				},
			},
		}},
	}

	if !assertEquals(a.Nodes[0], tree[0]) {
		t.Fatalf("Array bracked did not parse correctly")
	}
	if !assertEquals(a.Nodes[1], tree[1]) {
		t.Fatalf("Array bracked did not parse correctly")
	}
}
