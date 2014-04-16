package format

import (
	"io"
	"strings"

	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

func (f *formatWalker) Walk(node ast.Node) error {
	switch n := node.(type) {
	case ast.ExpressionStmt:
		f.Walk(n.Expression)
		f.printToken(token.StatementEnd)
		f.print("\n")
	case *ast.Variable:
		f.printToken(token.VariableOperator)
		f.Walk(n.Name)
	case ast.Variable:
		f.printToken(token.VariableOperator)
		f.Walk(n.Name)
	case ast.Identifier:
		f.print(n.Value)
	case *ast.Literal:
		f.print(n.Value)
	case ast.Literal:
		f.print(n.Value)
	case ast.EchoStmt:
		f.printEcho(&n)
	case *ast.EchoStmt:
		f.printEcho(n)
	case ast.IfStmt:
		f.printIf(&n)
	case *ast.IfStmt:
		f.printIf(n)
	case *ast.Block:
		for _, stmt := range n.Statements {
			f.Walk(stmt)
		}
	}
	return nil
}

func (f *formatWalker) print(s string) {
	io.WriteString(f.w, s)
}

func (f *formatWalker) printToken(t token.Token) {
	if s, ok := tokenMap[t]; ok {
		io.WriteString(f.w, s)
		return
	}
	io.WriteString(f.w, t.String())
}

func (f *formatWalker) printTab() {
	io.WriteString(f.w, strings.Repeat(f.Indent, f.tabLevel))
}

func (f *formatWalker) printIf(node *ast.IfStmt) error {
	f.printTab()
	f.printToken(token.If)
	f.print(" ")
	f.printToken(token.OpenParen)
	f.Walk(node.Condition)
	f.printToken(token.CloseParen)
	f.print(" ")
	f.printToken(token.BlockBegin)
	f.print("\n")
	f.tabLevel += 1
	f.Walk(node.TrueBranch)
	f.print("\n")
	f.tabLevel -= 1
	f.printToken(token.BlockEnd)
	if node.FalseBranch != nil {
		f.print(" ")
		f.printToken(token.Else)
		f.print(" ")
		f.printToken(token.BlockBegin)
		f.print("\n")
		f.tabLevel += 1
		f.Walk(node.FalseBranch)
		f.tabLevel -= 1
		f.print("\n")
		f.printToken(token.BlockEnd)
		f.print("\n")
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
