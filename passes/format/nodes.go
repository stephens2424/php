package format

import (
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/token"
)

func (f *formatWalker) printIf(node *ast.IfStmt) error {
	f.printTab()
	f.printf("if (%?) {", node.Condition)
	f.println()
	f.tabLevel += 1
	f.Walk(node.TrueBranch)
	f.tabLevel -= 1
	f.printTab()
	f.printToken(token.BlockEnd)
	if node.FalseBranch != nil && len(node.FalseBranch.Children()) != 0 {
		f.printf(" else {\n")
		f.tabLevel += 1
		f.Walk(node.FalseBranch)
		f.tabLevel -= 1
		f.printTab()
		f.printf("}\n")
	} else {
		f.println()
	}
	return nil
}

func (f *formatWalker) printEcho(node *ast.EchoStmt) error {
	f.printTab()
	f.printToken(token.Echo)
	for i, expr := range node.Expressions {
		f.print(" ")
		f.Walk(expr)
		if i != len(node.Expressions)-1 {
			f.printToken(token.Comma)
		}
	}
	f.printToken(token.StatementEnd)
	f.print("\n")
	return nil
}
