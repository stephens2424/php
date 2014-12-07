package printer

import (
	"fmt"
	"bytes"
	"github.com/stephens2424/php/ast"
)

type Printer struct {}

func (p *Printer) PrintNode(n ast.Node) string {
	switch n.(type) {
	default:
		return fmt.Sprintf(`/* Unsupported node type: %T */`, n)
	}
}

func (p *Printer) PrintIdentifier(i *ast.Identifier) string {
	return i.Value
}


func (p *Printer) PrintVariable(v *ast.Variable) string {
	return fmt.Sprintf("$%s", p.PrintNode(v.Name))
}

func (p *Printer) PrintGlobalDeclaration(g *ast.GlobalDeclaration) string {
	buf := bytes.NewBufferString("global ")
	for i, id := range g.Identifiers {
		buf.WriteString(p.PrintNode(id))
		if i+1 < len(g.Identifiers) {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func (p *Printer) PrintEmptyStatement(e *ast.EmptyStatement) string { return ";" }


func (p *Printer) PrintBinaryExpression(b *ast.BinaryExpression) string {
	return fmt.Sprintf("%s %s %s", b.Antecedent, b.Operator, b.Subsequent)
}


func (p *Printer) PrintTernaryExpression(t *ast.TernaryExpression) string {
	return fmt.Sprintf("%s ? %s : %s", t.Condition, t.True, t.False)
}
func (p *Printer) PrintUnaryExpression(u *ast.UnaryExpression) string {
	if u.Preceding {
		return fmt.Sprintf("%s%s", u.Operator, u.Operand)
	}
	return fmt.Sprintf("%s%s", u.Operand, u.Operator)
}
func (p *Printer) PrintEchoStmt(e *ast.EchoStmt) string {
	buf := bytes.NewBufferString("echo ")
	for i, expr := range e.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(p.PrintNode(expr))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintReturnStmt(r *ast.ReturnStmt) string {
	buf := bytes.NewBufferString("return")
	if r.Expression != nil {
		fmt.Fprintf(buf, " %s", p.PrintNode(r.Expression))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintBreakStmt(b *ast.BreakStmt) string {
	buf := bytes.NewBufferString("break")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", p.PrintNode(b.Expression))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintContinueStmt(b *ast.ContinueStmt) string {
	buf := bytes.NewBufferString("continue")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", p.PrintNode(b.Expression))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintThrowStmt(b *ast.ThrowStmt) string {
	buf := bytes.NewBufferString("throw")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", p.PrintNode(b.Expression))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintInclude(e *ast.Include) string {
	buf := bytes.NewBufferString("include ")
	for i, expr := range e.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(p.PrintNode(expr))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintExitStmt(b *ast.ExitStmt) string {
	buf := bytes.NewBufferString("exit")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", p.PrintNode(b.Expression))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintNewExpression(b *ast.NewExpression) string {
	buf := bytes.NewBufferString("new ")
	buf.WriteString(p.PrintNode(b.Class))
	buf.WriteString("(")
	for i, arg := range b.Arguments {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(p.PrintNode(arg))
	}
	buf.WriteString(")")
	return buf.String()
}
func (p *Printer) PrintAssignmentExpression(a *ast.AssignmentExpression) string {
	buf := bytes.NewBufferString(p.PrintNode(a.Assignee))
	buf.WriteString(" ")
	buf.WriteString(a.Operator)
	buf.WriteString(" ")
	buf.WriteString(p.PrintNode(a.Value))
	return buf.String()
}
func (p *Printer) PrintFunctionCallStmt(f *ast.FunctionCallStmt) string {
	buf := bytes.NewBufferString(p.PrintNode(f.FunctionCallExpression))
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintFunctionCallExpression(f *ast.FunctionCallExpression) string {
	buf := bytes.NewBufferString(p.PrintNode(f.FunctionName))
	buf.WriteString("(")
	for i, arg := range f.Arguments {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(p.PrintNode(arg))
	}
	buf.WriteString(")")
	return buf.String()
}


func (p *Printer) PrintBlock(b *ast.Block) string {
	buf := &bytes.Buffer{}
	for _, s := range b.Statements {
		buf.WriteString(p.PrintNode(s))
	}
	return buf.String()
}
func (p *Printer) PrintFunctionStmt(f *ast.FunctionStmt) string {
	return fmt.Sprintf("%s%s", p.PrintNode(f.FunctionDefinition), p.PrintNode(f.Body))
}
func (p *Printer) PrintAnonymousFunction(a *ast.AnonymousFunction) string {
	buf := bytes.NewBufferString("function (")
	for i, arg := range a.Arguments {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(p.PrintNode(arg))
	}
	buf.WriteString(")")
	if len(a.ClosureVariables) > 0 {
		fmt.Fprint(buf, " use (")
		for i, arg := range a.ClosureVariables {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(p.PrintNode(arg))
		}
		buf.WriteString(")")
	}
	buf.WriteString(p.PrintNode(a.Body))
	return buf.String()
}

func (p *Printer) PrintFunctionDefinition(fd *ast.FunctionDefinition) string {
	buf := bytes.NewBufferString("function ")
	buf.WriteString(fd.Name)
	buf.WriteString(" (")
	for i, arg := range fd.Arguments {
		buf.WriteString(p.PrintNode(arg))
		if i+1 < len(fd.Arguments) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	return buf.String()
}
func (p *Printer) PrintFunctionArgument(fa *ast.FunctionArgument) string {
	buf := &bytes.Buffer{}
	if fa.TypeHint != "" {
		fmt.Fprint(buf, fa.TypeHint, "")
	}
	buf.WriteString(p.PrintNode(fa.Variable))
	if fa.Default != nil {
		fmt.Fprint(buf, " =", p.PrintNode(fa.Default))
	}
	return buf.String()
}
func (p *Printer) PrintClass(c *ast.Class) string {
	buf := bytes.NewBufferString("class ")
	buf.WriteString(c.Name)
	if c.Extends != "" {
		fmt.Fprintf(buf, " extends %s", c.Extends)
	}
	for i, imp := range c.Implements {
		if i > 0 {
			buf.WriteString(",")
		} else {
			buf.WriteString("implements ")
		}
		buf.WriteString(imp)
	}
	buf.WriteString(" {\n")
	for _, c := range c.Constants {
		buf.WriteString(p.PrintNode(c))
	}
	for _, pr := range c.Properties {
		buf.WriteString(p.PrintNode(pr))
	}
	for _, m := range c.Methods {
		buf.WriteString(p.PrintNode(m))
	}
	buf.WriteString("}")
	return buf.String()
}

func (p *Printer) PrintInterface(i *ast.Interface) string {
	buf := bytes.NewBufferString("interface ")
	buf.WriteString(i.Name)

	for i, imp := range i.Inherits {
		if i > 0 {
			buf.WriteString(", ")
		} else {
			buf.WriteString("implements ")
		}
		buf.WriteString(imp)
	}

	buf.WriteString(" {")
	for _, c := range i.Constants {
		buf.WriteString(p.PrintNode(c))
	}

	for _, m := range i.Methods {
		buf.WriteString(p.PrintNode(m))
	}

	buf.WriteString("}")
	return buf.String()
}
func (p *Printer) PrintProperty(pr *ast.Property) string {
	buf := &bytes.Buffer{}
	buf.WriteString(pr.Visibility.Token().String())
	fmt.Fprintf(buf, " %s", pr.Name)
	if pr.Initialization != nil {
		fmt.Fprintf(buf, " = %s", p.PrintNode(pr.Initialization))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p *Printer) PrintPropertyExpression(pr *ast.PropertyExpression) string {
	return fmt.Sprintf("%s->%s", p.PrintNode(pr.Receiver), p.PrintNode(pr.Name))
}

func (p *Printer) PrintClassExpression(c *ast.ClassExpression) string {
	return fmt.Sprintf("%s::%s", p.PrintNode(c.Receiver), p.PrintNode(c.Expression))
}
func (p *Printer) PrintMethod(m *ast.Method) string {
	return fmt.Sprintf("%s %s", m.Visibility.Token().String(), p.PrintNode(m.FunctionStmt))
}
func (p *Printer) PrintMethodCallExpression(m *ast.MethodCallExpression) string {
	return fmt.Sprintf("%s->%s", p.PrintNode(m.Receiver), p.PrintNode(m.FunctionCallExpression))
}
func (p *Printer) PrintIfStmt(i *ast.IfStmt) string {
	str := fmt.Sprintf("if (%s) {\n%s\n}", i.Condition, i.TrueBranch)
	if i.FalseBranch != nil {
		str += fmt.Sprintf(" else {\n%s\n}", i.FalseBranch)
	}
	return str
}

func (p *Printer) PrintSwitchStmt(s *ast.SwitchStmt) string {
	str := fmt.Sprintf("switch (%s) {\n", s.Expression)
	for _, c := range s.Cases {
		str += fmt.Sprintf("%s\n", p.PrintNode(c))
	}
	if s.DefaultCase != nil {
		str += fmt.Sprintf("default:\n%s", p.PrintNode(s.DefaultCase))
	}
	str += "}"
	return str
}
func (p *Printer) PrintSwitchCase(s *ast.SwitchCase) string {
	return fmt.Sprintf("case %s:\n%s\n", s.Expression, s.Block)
}
func (p *Printer) PrintForStmt(f *ast.ForStmt) string {
	str := fmt.Sprintf("for (")
	for i, e := range f.Initialization {
		if i > 0 {
			str += ", "
		}
		str += p.PrintNode(e)
	}
	for i, e := range f.Termination {
		if i > 0 {
			str += ", "
		}
		str += p.PrintNode(e)
	}
	for i, e := range f.Iteration {
		if i > 0 {
			str += ", "
		}
		str += p.PrintNode(e)
	}
			str += ") "
			str += p.PrintNode(f.LoopBlock)
			return str
}
func (p *Printer) PrintWhileStmt(w *ast.WhileStmt) string {
	return fmt.Sprintf("while (%s) %s", w.Termination, w.LoopBlock)
}
func (p *Printer) PrintDoWhileStmt(w *ast.DoWhileStmt) string {
	return fmt.Sprintf("do %s while (%s);", w.LoopBlock, w.Termination)
}
func (p *Printer) PrintTryStmt(t *ast.TryStmt) string {
	str := fmt.Sprintf("try %s", p.PrintNode(t.TryBlock))
	for _, c := range t.CatchStmts {
		str += p.PrintNode(c)
	}
	if t.FinallyBlock != nil {
		str += fmt.Sprintf("finally %s", p.PrintNode(t.FinallyBlock))
	}
	return str
}
func (p *Printer) PrintCatchStmt(c *ast.CatchStmt) string {
	return fmt.Sprintf("catch (%s %s) %s", c.CatchType, p.PrintNode(c.CatchVar), p.PrintNode(c.CatchBlock))
}
func (p *Printer) PrintLiteral(l *ast.Literal) string {
	switch l.Type {
	case ast.String:
		return  l.Value
	case ast.Integer, ast.Float:
		return  l.Value
	case ast.Boolean:
		return  l.Value
	case ast.Null:
		return  "null"
	}
	panic("invalid literal type")
}
func (p *Printer) PrintForeachStmt(f *ast.ForeachStmt) string {
	str := fmt.Sprintf("foreach (%s as ", f.Source)
	if f.Key != nil {
		str += fmt.Sprintf("%s => ", f.Key)
	}
	str += fmt.Sprintf("%s) %s", p.PrintNode(f.Value), p.PrintNode(f.LoopBlock))
	return str
}
func (p *Printer) PrintArrayExpression(a *ast.ArrayExpression) string {
	str := fmt.Sprintf("array(")
	for i, pair := range a.Pairs {
		if i > 0 {
			str += ", "
		}
		str += p.PrintNode(pair)
	}
	str += ")"
	return str
}
func (p *Printer) PrintArrayPair(pr *ast.ArrayPair) string {
	if pr.Key != nil {
	return fmt.Sprintf("%s => %s", p.PrintNode(pr.Key), p.PrintNode(pr.Value))
	}
	return fmt.Sprintf("%s", p.PrintNode(pr.Value))
}
func (p *Printer) PrintArrayLookupExpression(a *ast.ArrayLookupExpression) string {
	return fmt.Sprintf("%s[%s]", a.Array, p.PrintNode(a.Index))
}
func (p *Printer) PrintArrayAppendExpression(a *ast.ArrayAppendExpression) string {
	return fmt.Sprintf("%s[]", a.Array)
}
func (p *Printer) PrintShellCommand(s *ast.ShellCommand) string {
	return  s.Command
}
func (p *Printer) PrintListStatement(l *ast.ListStatement) string {
	str := fmt.Sprintf("list(")
	for i, a := range l.Assignees {
		if i > 0 {
			str += ", "
		}
		str += p.PrintNode(a)
	}
	str += fmt.Sprintf(") = ")
	str += p.PrintNode(l.Value)
	return str
}

func (p *Printer) PrintStaticVariableDeclaration(s *ast.StaticVariableDeclaration) string {
	str := fmt.Sprintf("static ")
	for i, d := range s.Declarations {
		if i > 0 {
			str += ", "
		}
		str += p.PrintNode(d)
	}
	str += ";\n"
	return str
}
func (p *Printer) PrintDeclareBlock(d *ast.DeclareBlock) string {
	buf := bytes.NewBufferString("declare (")
	for i, decl := range d.Declarations {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(decl)
	}
	buf.WriteString(") {")
	buf.WriteString(p.PrintNode(d.Statements))
	buf.WriteString("}")
	return buf.String()
}
