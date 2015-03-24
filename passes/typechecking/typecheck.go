package typecheck

import (
	"fmt"

	"github.com/stephens2424/php/ast"
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
		for _, branch := range n.Branches {
			fmt.Println("parsed condition:", branch.Condition)
			if !branch.Condition.EvaluatesTo().Contains(ast.Boolean) {
				w.Errorf("If condition does not evaluate to boolean")
			}
		}
	}
}
