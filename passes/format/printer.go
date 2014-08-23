package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/token"
)

func (f *formatWalker) Walk(node ast.Node) error {
	f.printNode(node)
	return nil
}

func (f *formatWalker) printNode(node ast.Node) {
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
	case ast.Class:
		f.printToken(token.Class)
		f.printf(" %s {", n.Name)
		f.tabLevel += 1
		for _, child := range n.Children() {
			f.Walk(child)
		}
	case ast.Property:
		f.printTabbedLine()
		f.printVisibility(n.Visibility)
		f.printf(" %s", n.Name)
		if n.Initialization != nil {
			f.print(" = ")
			f.printNode(n.Initialization)
		}
		f.print(";")
	default:
		f.printTabbedLine()
		f.printf("// unimplemented %T\n", n)
		for _, child := range n.Children() {
			f.Walk(child)
		}
	}
}

func (f *formatWalker) printf(fmtString string, i ...interface{}) {
	num := 0
	inPercent := false
	for _, r := range fmtString {
		if inPercent && r == '?' {
			switch n := i[num].(type) {
			case token.Token:
				f.printToken(n)
				num += 1
			case ast.Node:
				f.printNode(n)
				num += 1
			default:
				fmt.Fprint(f.w, n)
			}
		} else if inPercent && r == '%' {
			io.WriteString(f.w, "%")
		} else if inPercent {
			fmt.Fprintf(f.w, "%"+string(byte(r)), i[num])
		} else if r == '%' {
			inPercent = true
			continue
		} else {
			fmt.Fprint(f.w, string(byte(r)))
		}
		inPercent = false
	}
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

func (f *formatWalker) printTabbedLine() {
	f.println()
	f.printTab()
}

func (f *formatWalker) println() {
	io.WriteString(f.w, "\n")
}

var visibilityMap = map[ast.Visibility]string{
	ast.Private:   "private",
	ast.Protected: "protected",
	ast.Public:    "public",
}

func (f *formatWalker) printVisibility(v ast.Visibility) {
	f.print(visibilityMap[v])
}
