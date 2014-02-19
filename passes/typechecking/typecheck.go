package typecheck

import (
	"fmt"
	"stephensearles.com/php/ast"
)

type Walker struct {
	ast.DefaultWalker
}

func (w *Walker) Walk(node ast.Node) {
	switch n := node.(type) {
	case ast.Block:
		for _, stmt := range n.Statements {
			w.Walk(stmt)
		}
	case *ast.IfStmt:
		fmt.Println("parsed condition:", n.Condition)
		if !n.Condition.EvaluatesTo().Contains(ast.Boolean) {
			w.Errorf("If condition does not evaluate to boolean")
		}
	}
}
