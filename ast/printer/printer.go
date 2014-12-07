package printer

import (
	"bytes"
	"fmt"
	"github.com/stephens2424/php/ast"
	"io"
)

type Printer struct{}

func (p *Printer) PrintNode(w io.Writer, node ast.Node) {
	switch n := node.(type) {
	case *ast.AnonymousFunction:
		p.PrintAnonymousFunction(w, n)
	case *ast.ArrayAppendExpression:
		p.PrintArrayAppendExpression(w, n)
	case *ast.ArrayExpression:
		p.PrintArrayExpression(w, n)
	case *ast.ArrayLookupExpression:
		p.PrintArrayLookupExpression(w, n)
	case *ast.ArrayPair:
		p.PrintArrayPair(w, n)
	case *ast.AssignmentExpression:
		p.PrintAssignmentExpression(w, n)
	case *ast.BinaryExpression:
		p.PrintBinaryExpression(w, n)
	case *ast.Block:
		p.PrintBlock(w, n)
	case *ast.BreakStmt:
		p.PrintBreakStmt(w, n)
	case *ast.CatchStmt:
		p.PrintCatchStmt(w, n)
	case *ast.Class:
		p.PrintClass(w, n)
	case *ast.ClassExpression:
		p.PrintClassExpression(w, n)
	case *ast.Constant:
		p.PrintConstant(w, n)
	case *ast.ConstantExpression:
		p.PrintConstantExpression(w, n)
	case *ast.ContinueStmt:
		p.PrintContinueStmt(w, n)
	case *ast.DeclareBlock:
		p.PrintDeclareBlock(w, n)
	case *ast.DoWhileStmt:
		p.PrintDoWhileStmt(w, n)
	case *ast.EchoStmt:
		p.PrintEchoStmt(w, n)
	case *ast.EmptyStatement:
		p.PrintEmptyStatement(w, n)
	case *ast.ExitStmt:
		p.PrintExitStmt(w, n)
	case *ast.ExpressionStmt:
		p.PrintExpressionStmt(w, n)
	case *ast.ForStmt:
		p.PrintForStmt(w, n)
	case *ast.ForeachStmt:
		p.PrintForeachStmt(w, n)
	case *ast.FunctionArgument:
		p.PrintFunctionArgument(w, n)
	case *ast.FunctionCallExpression:
		p.PrintFunctionCallExpression(w, n)
	case *ast.FunctionCallStmt:
		p.PrintFunctionCallStmt(w, n)
	case *ast.FunctionDefinition:
		p.PrintFunctionDefinition(w, n)
	case *ast.FunctionStmt:
		p.PrintFunctionStmt(w, n)
	case *ast.GlobalDeclaration:
		p.PrintGlobalDeclaration(w, n)
	case *ast.Identifier:
		p.PrintIdentifier(w, n)
	case *ast.IfStmt:
		p.PrintIfStmt(w, n)
	case *ast.Include:
		p.PrintInclude(w, n)
	case *ast.IncludeStmt:
		p.PrintIncludeStmt(w, n)
	case *ast.Interface:
		p.PrintInterface(w, n)
	case *ast.ListStatement:
		p.PrintListStatement(w, n)
	case *ast.Literal:
		p.PrintLiteral(w, n)
	case *ast.Method:
		p.PrintMethod(w, n)
	case *ast.MethodCallExpression:
		p.PrintMethodCallExpression(w, n)
	case *ast.NewExpression:
		p.PrintNewExpression(w, n)
	case *ast.Property:
		p.PrintProperty(w, n)
	case *ast.PropertyExpression:
		p.PrintPropertyExpression(w, n)
	case *ast.ReturnStmt:
		p.PrintReturnStmt(w, n)
	case *ast.ShellCommand:
		p.PrintShellCommand(w, n)
	case *ast.StaticVariableDeclaration:
		p.PrintStaticVariableDeclaration(w, n)
	case *ast.SwitchCase:
		p.PrintSwitchCase(w, n)
	case *ast.SwitchStmt:
		p.PrintSwitchStmt(w, n)
	case *ast.TernaryExpression:
		p.PrintTernaryExpression(w, n)
	case *ast.ThrowStmt:
		p.PrintThrowStmt(w, n)
	case *ast.TryStmt:
		p.PrintTryStmt(w, n)
	case *ast.UnaryExpression:
		p.PrintUnaryExpression(w, n)
	case *ast.Variable:
		p.PrintVariable(w, n)
	case *ast.WhileStmt:
		p.PrintWhileStmt(w, n)
	default:
		fmt.Fprintf(w, `/* Unsupported node type: %T */`, n)
	}
}

func (p *Printer) PrintIdentifier(w io.Writer, i *ast.Identifier) {
	io.WriteString(w, i.Value)
}

func (p *Printer) PrintVariable(w io.Writer, v *ast.Variable) {
	io.WriteString(w, "$")
	p.PrintNode(w, v.Name)
}

func (p *Printer) PrintGlobalDeclaration(w io.Writer, g *ast.GlobalDeclaration) {
	io.WriteString(w, "global ")
	for i, id := range g.Identifiers {
		p.PrintNode(w, id)
		if i+1 < len(g.Identifiers) {
			io.WriteString(w, ", ")
		}
	}
}

func (p *Printer) PrintEmptyStatement(w io.Writer, e *ast.EmptyStatement) {}

func (p *Printer) PrintBinaryExpression(w io.Writer, b *ast.BinaryExpression) {
	fmt.Fprintf(w, "%s %s %s", b.Antecedent, b.Operator, b.Subsequent)
}

func (p *Printer) PrintTernaryExpression(w io.Writer, t *ast.TernaryExpression) {
	fmt.Fprintf(w, "%s ? %s : %s", t.Condition, t.True, t.False)
}
func (p *Printer) PrintUnaryExpression(w io.Writer, u *ast.UnaryExpression) {
	if u.Preceding {
		fmt.Fprintf(w, "%s%s", u.Operator, u.Operand)
	}
	fmt.Fprintf(w, "%s%s", u.Operand, u.Operator)
}

func (p *Printer) PrintEchoStmt(w io.Writer, e *ast.EchoStmt) {
	io.WriteString(w, "echo ")
	for i, expr := range e.Expressions {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, expr)
	}
	io.WriteString(w, ";")
}

func (p *Printer) PrintReturnStmt(w io.Writer, r *ast.ReturnStmt) {
	io.WriteString(w, "return")
	if r.Expression != nil {
		p.PrintNode(w, r.Expression)
	}
	io.WriteString(w, ";")
}
func (p *Printer) PrintBreakStmt(w io.Writer, b *ast.BreakStmt) {
	io.WriteString(w, "break")
	if b.Expression != nil {
		p.PrintNode(w, b.Expression)
	}
	io.WriteString(w, ";")

}
func (p *Printer) PrintContinueStmt(w io.Writer, b *ast.ContinueStmt) {
	io.WriteString(w, "continue")
	if b.Expression != nil {
		p.PrintNode(w, b.Expression)
	}
	io.WriteString(w, ";")

}
func (p *Printer) PrintThrowStmt(w io.Writer, b *ast.ThrowStmt) {
	io.WriteString(w, "throw")
	if b.Expression != nil {
		p.PrintNode(w, b.Expression)
	}
	io.WriteString(w, ";")

}
func (p *Printer) PrintInclude(w io.Writer, e *ast.Include) {
	io.WriteString(w, "include ")
	for i, expr := range e.Expressions {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, expr)
	}
	io.WriteString(w, ";")

}
func (p *Printer) PrintExitStmt(w io.Writer, b *ast.ExitStmt) {
	io.WriteString(w, "exit")
	if b.Expression != nil {
		p.PrintNode(w, b.Expression)
	}
	io.WriteString(w, ";")

}
func (p *Printer) PrintNewExpression(w io.Writer, b *ast.NewExpression) {
	io.WriteString(w, "new ")
	p.PrintNode(w, b.Class)
	io.WriteString(w, "(")
	for i, arg := range b.Arguments {
		if i > 0 {
			io.WriteString(w, ",")
		}
		p.PrintNode(w, arg)
	}
	io.WriteString(w, ")")

}
func (p *Printer) PrintAssignmentExpression(w io.Writer, a *ast.AssignmentExpression) {
	p.PrintNode(w, a.Assignee)
	io.WriteString(w, " ")
	io.WriteString(w, a.Operator)
	io.WriteString(w, " ")
	p.PrintNode(w, a.Value)
}

func (p *Printer) PrintFunctionCallStmt(w io.Writer, f *ast.FunctionCallStmt) {
	p.PrintNode(w, f.FunctionCallExpression)
	io.WriteString(w, ";")

}
func (p *Printer) PrintFunctionCallExpression(w io.Writer, f *ast.FunctionCallExpression) {
	p.PrintNode(w, f.FunctionName)
	io.WriteString(w, "(")
	for i, arg := range f.Arguments {
		if i > 0 {
			io.WriteString(w, ",")
		}
		p.PrintNode(w, arg)
	}
	io.WriteString(w, ")")

}

func (p *Printer) PrintBlock(w io.Writer, b *ast.Block) {
	for _, s := range b.Statements {
		p.PrintNode(w, s)
		io.WriteString(w, "\n")
	}

}
func (p *Printer) PrintFunctionStmt(w io.Writer, f *ast.FunctionStmt) {
	p.PrintNode(w, f.FunctionDefinition)
	p.PrintNode(w, f.Body)
}
func (p *Printer) PrintAnonymousFunction(w io.Writer, a *ast.AnonymousFunction) {
	io.WriteString(w, "function (")
	for i, arg := range a.Arguments {
		if i > 0 {
			io.WriteString(w, ",")
		}
		p.PrintNode(w, arg)
	}
	io.WriteString(w, ")")
	if len(a.ClosureVariables) > 0 {
		fmt.Fprint(w, " use (")
		for i, arg := range a.ClosureVariables {
			if i > 0 {
				io.WriteString(w, ",")
			}
			p.PrintNode(w, arg)
		}
		io.WriteString(w, ")")
	}
	p.PrintNode(w, a.Body)

}

func (p *Printer) PrintFunctionDefinition(w io.Writer, fd *ast.FunctionDefinition) {
	io.WriteString(w, "function ")
	io.WriteString(w, fd.Name)
	io.WriteString(w, " (")
	for i, arg := range fd.Arguments {
		p.PrintNode(w, arg)
		if i+1 < len(fd.Arguments) {
			io.WriteString(w, ",")
		}
	}
	io.WriteString(w, ")")

}
func (p *Printer) PrintFunctionArgument(w io.Writer, fa *ast.FunctionArgument) {
	buf := &bytes.Buffer{}
	if fa.TypeHint != "" {
		fmt.Fprint(buf, fa.TypeHint, "")
	}
	p.PrintNode(w, fa.Variable)
	if fa.Default != nil {
		io.WriteString(w, " =")
		p.PrintNode(w, fa.Default)
	}

}
func (p *Printer) PrintClass(w io.Writer, c *ast.Class) {
	io.WriteString(w, "class ")
	io.WriteString(w, c.Name)
	if c.Extends != "" {
		fmt.Fprintf(w, " extends %s", c.Extends)
	}
	for i, imp := range c.Implements {
		if i > 0 {
			io.WriteString(w, ",")
		} else {
			io.WriteString(w, "implements ")
		}
		io.WriteString(w, imp)
	}
	io.WriteString(w, " {\n")
	for _, c := range c.Constants {
		p.PrintNode(w, c)
	}
	for _, pr := range c.Properties {
		p.PrintNode(w, pr)
	}
	for _, m := range c.Methods {
		p.PrintNode(w, m)
	}
	io.WriteString(w, "}")

}

func (p *Printer) PrintInterface(w io.Writer, i *ast.Interface) {
	io.WriteString(w, "interface ")
	io.WriteString(w, i.Name)

	for i, imp := range i.Inherits {
		if i > 0 {
			io.WriteString(w, ", ")
		} else {
			io.WriteString(w, "implements ")
		}
		io.WriteString(w, imp)
	}

	io.WriteString(w, " {")
	for _, c := range i.Constants {
		p.PrintNode(w, c)
	}

	for _, m := range i.Methods {
		p.PrintNode(w, m)
	}

	io.WriteString(w, "}")

}
func (p *Printer) PrintProperty(w io.Writer, pr *ast.Property) {
	buf := &bytes.Buffer{}
	io.WriteString(w, pr.Visibility.Token().String())
	fmt.Fprintf(buf, " %s", pr.Name)
	if pr.Initialization != nil {
		p.PrintNode(w, pr.Initialization)
	}
	io.WriteString(w, ";")

}
func (p *Printer) PrintPropertyExpression(w io.Writer, pr *ast.PropertyExpression) {
	p.PrintNode(w, pr.Receiver)
	io.WriteString(w, "->")
	p.PrintNode(w, pr.Name)
}

func (p *Printer) PrintClassExpression(w io.Writer, c *ast.ClassExpression) {
	p.PrintNode(w, c.Receiver)
	io.WriteString(w, "::")
	p.PrintNode(w, c.Expression)
}
func (p *Printer) PrintMethod(w io.Writer, m *ast.Method) {
	fmt.Fprintf(w, "%s ", m.Visibility.Token().String())
	p.PrintNode(w, m.FunctionStmt)
}
func (p *Printer) PrintMethodCallExpression(w io.Writer, m *ast.MethodCallExpression) {
	p.PrintNode(w, m.Receiver)
	io.WriteString(w, "->")
	p.PrintNode(w, m.FunctionCallExpression)
}
func (p *Printer) PrintIfStmt(w io.Writer, i *ast.IfStmt) {
	fmt.Fprintf(w, "if (%s) {\n%s\n}", i.Condition, i.TrueBranch)
	if i.FalseBranch != nil {
		fmt.Fprintf(w, " else {\n%s\n}", i.FalseBranch)
	}

}

func (p *Printer) PrintSwitchStmt(w io.Writer, s *ast.SwitchStmt) {
	fmt.Fprintf(w, "switch (%s) {\n", s.Expression)
	for _, c := range s.Cases {
		p.PrintNode(w, c)
		io.WriteString(w, "\n")
	}
	if s.DefaultCase != nil {
		fmt.Fprintf(w, "default:\n")
		p.PrintNode(w, s.DefaultCase)
	}
	io.WriteString(w, "}")
}

func (p *Printer) PrintSwitchCase(w io.Writer, s *ast.SwitchCase) {
	io.WriteString(w, "case ")
	p.PrintNode(w, s.Expression)
	io.WriteString(w, ":\n")
	p.PrintNode(w, s.Block)
	io.WriteString(w, "\n")
}
func (p *Printer) PrintForStmt(w io.Writer, f *ast.ForStmt) {
	fmt.Fprintf(w, "for (")
	for i, e := range f.Initialization {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, e)
	}
	for i, e := range f.Termination {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, e)
	}
	for i, e := range f.Iteration {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, e)
	}
	io.WriteString(w, ") ")
	p.PrintNode(w, f.LoopBlock)

}
func (p *Printer) PrintWhileStmt(w io.Writer, wh *ast.WhileStmt) {
	fmt.Fprintf(w, "while (%s) %s", wh.Termination, wh.LoopBlock)
}
func (p *Printer) PrintDoWhileStmt(w io.Writer, wh *ast.DoWhileStmt) {
	fmt.Fprintf(w, "do %s while (%s);", wh.LoopBlock, wh.Termination)
}
func (p *Printer) PrintTryStmt(w io.Writer, t *ast.TryStmt) {
	fmt.Fprintf(w, "try ")
	p.PrintNode(w, t.TryBlock)
	for _, c := range t.CatchStmts {
		p.PrintNode(w, c)
	}
	if t.FinallyBlock != nil {
		fmt.Fprintf(w, "finally ")
		p.PrintNode(w, t.FinallyBlock)
	}

}
func (p *Printer) PrintCatchStmt(w io.Writer, c *ast.CatchStmt) {
	fmt.Fprintf(w, "catch (%s ", c.CatchType)
	p.PrintNode(w, c.CatchVar)
	io.WriteString(w, ") ")
	p.PrintNode(w, c.CatchBlock)
}

func (p *Printer) PrintLiteral(w io.Writer, l *ast.Literal) {
	switch l.Type {
	case ast.String:
		io.WriteString(w, l.Value)
	case ast.Integer, ast.Float:
		io.WriteString(w, l.Value)
	case ast.Boolean:
		io.WriteString(w, l.Value)
	case ast.Null:
		io.WriteString(w, "null")
	}
	panic("invalid literal type")
}

func (p *Printer) PrintForeachStmt(w io.Writer, f *ast.ForeachStmt) {
	fmt.Fprintf(w, "foreach (%s as ", f.Source)
	if f.Key != nil {
		fmt.Fprintf(w, "%s => ", f.Key)
	}
	p.PrintNode(w, f.Value)
	io.WriteString(w, ") ")
	p.PrintNode(w, f.LoopBlock)
}

func (p *Printer) PrintArrayExpression(w io.Writer, a *ast.ArrayExpression) {
	fmt.Fprintf(w, "array(")
	for i, pair := range a.Pairs {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, pair)
	}
	io.WriteString(w, ")")
}

func (p *Printer) PrintArrayPair(w io.Writer, pr *ast.ArrayPair) {
	if pr.Key != nil {
		p.PrintNode(w, pr.Key)
		fmt.Fprintf(w, " => ")
		p.PrintNode(w, pr.Value)
	}
	p.PrintNode(w, pr.Value)
}

func (p *Printer) PrintArrayLookupExpression(w io.Writer, a *ast.ArrayLookupExpression) {
	p.PrintNode(w, a.Array)
	io.WriteString(w, "[")
	p.PrintNode(w, a.Index)
	io.WriteString(w, "]")
}

func (p *Printer) PrintArrayAppendExpression(w io.Writer, a *ast.ArrayAppendExpression) {
	fmt.Fprintf(w, "%s[]", a.Array)
}

func (p *Printer) PrintShellCommand(w io.Writer, s *ast.ShellCommand) {
	io.WriteString(w, s.Command)
}

func (p *Printer) PrintListStatement(w io.Writer, l *ast.ListStatement) {
	fmt.Fprintf(w, "list(")
	for i, a := range l.Assignees {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, a)
	}
	io.WriteString(w, ") =")
	p.PrintNode(w, l.Value)
}

func (p *Printer) PrintStaticVariableDeclaration(w io.Writer, s *ast.StaticVariableDeclaration) {
	fmt.Fprintf(w, "static ")
	for i, d := range s.Declarations {
		if i > 0 {
			io.WriteString(w, ", ")
		}
		p.PrintNode(w, d)
	}
	io.WriteString(w, ";\n")
}
func (p *Printer) PrintDeclareBlock(w io.Writer, d *ast.DeclareBlock) {
	io.WriteString(w, "declare (")
	for i, decl := range d.Declarations {
		if i > 0 {
			io.WriteString(w, ",")
		}
		io.WriteString(w, decl)
	}
	io.WriteString(w, ") {")
	p.PrintNode(w, d.Statements)
	io.WriteString(w, "}")
}

func (p *Printer) PrintConstant(w io.Writer, c *ast.Constant) {
	p.PrintNode(w, c.Variable.Name)
}

func (p *Printer) PrintConstantExpression(w io.Writer, c *ast.ConstantExpression) {
	p.PrintNode(w, c.Variable.Name)
}

func (p *Printer) PrintExpressionStmt(w io.Writer, c *ast.ExpressionStmt) {
	p.PrintNode(w, c.Expression)
	io.WriteString(w, ";")
}

func (p *Printer) PrintIncludeStmt(w io.Writer, c *ast.IncludeStmt) {
	p.PrintInclude(w, &c.Include)
	io.WriteString(w, ";")
}
