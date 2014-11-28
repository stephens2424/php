package main

import "github.com/stephens2424/php/ast"

var n ast.Node



func (i Identifier) Print(f Format) string {
	return i.Value
}


func (v Variable) Print(f Format) string {
	return fmt.Sprintf("$%s", v.Name.Print(f))
}
func (g GlobalDeclaration) Print(f Format) string {
	buf := bytes.NewBufferString("global ")
	for i, id := range g.Identifiers {
		buf.WriteString(id.Print(f))
		if i+1 < len(g.Identifiers) {
			buf.WriteString(", ")
		}
	}
	return buf.String()
}

func (e EmptyStatement) String() string        { return "" }
func (e EmptyStatement) Children() []Node      { return nil }
func (e EmptyStatement) Print(f Format) string { return ";" }


func (b BinaryExpression) Print(f Format) string {
	return fmt.Sprintf("%s %s %s", b.Antecedent, b.Operator, b.Subsequent)
}


func (t TernaryExpression) Print(f Format) string {
	return fmt.Sprintf("%s ? %s : %s", t.Condition, t.True, t.False)
}
func (u UnaryExpression) Print(f Format) string {
	if u.Preceding {
		return fmt.Sprintf("%s%s", u.Operator, u.Operand)
	}
	return fmt.Sprintf("%s%s", u.Operand, u.Operator)
}
func (e EchoStmt) Print(f Format) string {
	buf := bytes.NewBufferString("echo ")
	for i, expr := range e.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(expr.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (r ReturnStmt) Print(f Format) string {
	buf := bytes.NewBufferString("return")
	if r.Expression != nil {
		fmt.Fprintf(buf, " %s", r.Expression.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (b BreakStmt) Print(f Format) string {
	buf := bytes.NewBufferString("break")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", b.Expression.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (b ContinueStmt) Print(f Format) string {
	buf := bytes.NewBufferString("continue")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", b.Expression.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (b ThrowStmt) Print(f Format) string {
	buf := bytes.NewBufferString("throw")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", b.Expression.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (e Include) Print(f Format) string {
	buf := bytes.NewBufferString("include ")
	for i, expr := range e.Expressions {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(expr.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (b ExitStmt) Print(f Format) string {
	buf := bytes.NewBufferString("exit")
	if b.Expression != nil {
		fmt.Fprintf(buf, " %s", b.Expression.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (b NewExpression) Print(f Format) string {
	buf := bytes.NewBufferString("new ")
	buf.WriteString(b.Class.Print(f))
	buf.WriteString("(")
	for i, arg := range b.Arguments {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(arg.Print(f))
	}
	buf.WriteString(")")
	return buf.String()
}
func (a AssignmentExpression) Print(f Format) string {
	buf := bytes.NewBufferString(a.Assignee.Print(f))
	buf.WriteString(" ")
	buf.WriteString(a.Operator)
	buf.WriteString(" ")
	buf.WriteString(a.Value.Print(f))
	return buf.String()
}
func (f FunctionCallStmt) Print(fm Format) string {
	buf := bytes.NewBufferString(f.FunctionCallExpression.Print(fm))
	buf.WriteString(";")
	return buf.String()
}
func (f FunctionCallExpression) Print(fm Format) string {
	buf := bytes.NewBufferString(f.FunctionName.Print(fm))
	buf.WriteString("(")
	for i, arg := range f.Arguments {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(arg.Print(fm))
	}
	buf.WriteString(")")
	return buf.String()
}


func (b Block) Print(f Format) string {
	buf := &bytes.Buffer{}
	for _, s := range b.Statements {
		buf.WriteString(s.Print(f))
	}
	return buf.String()
}
func (f FunctionStmt) Print(fm Format) string {
	return fmt.Sprintf("%s%s", f.FunctionDefinition.Print(fm), f.Body.Print(fm))
}
func (a AnonymousFunction) Print(f Format) string {
	buf := bytes.NewBufferString("function (")
	for i, arg := range a.Arguments {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(arg.Print(f))
	}
	buf.WriteString(")")
	if len(a.ClosureVariables) > 0 {
		fmt.Fprint(buf, " use (")
		for i, arg := range a.ClosureVariables {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(arg.Print(f))
		}
		buf.WriteString(")")
	}
	buf.WriteString(a.Body.Print(f))
	return buf.String()
}

func (fd FunctionDefinition) Print(f Format) string {
	buf := bytes.NewBufferString("function ")
	buf.WriteString(fd.Name)
	buf.WriteString(" (")
	for i, arg := range fd.Arguments {
		buf.WriteString(arg.Print(f))
		if i+1 < len(fd.Arguments) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	return buf.String()
}
func (fa FunctionArgument) Print(f Format) string {
	buf := &bytes.Buffer{}
	if fa.TypeHint != "" {
		fmt.Fprint(buf, fa.TypeHint, "")
	}
	buf.WriteString(fa.Variable.Print(f))
	if fa.Default != nil {
		fmt.Fprint(buf, " =", fa.Default.Print(f))
	}
	return buf.String()
}
func (c Class) Print(f Format) string {
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
		buf.WriteString(c.Print(f))
	}
	for _, p := range c.Properties {
		buf.WriteString(p.Print(f))
	}
	for _, m := range c.Methods {
		buf.WriteString(m.Print(f))
	}
	buf.WriteString("}")
	return buf.String()
}

func (i Interface) Print(f Format) string {
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
		buf.WriteString(c.Print(f))
	}

	for _, m := range i.Methods {
		buf.WriteString(m.Print(f))
	}

	buf.WriteString("}")
	return buf.String()
}
func (p Property) Print(f Format) string {
	buf := &bytes.Buffer{}
	buf.WriteString(p.Visibility.Token().String())
	fmt.Fprintf(buf, " %s", p.Name)
	if p.Initialization != nil {
		fmt.Fprintf(buf, " = %s", p.Initialization.Print(f))
	}
	buf.WriteString(";")
	return buf.String()
}
func (p PropertyExpression) Print(f Format) string {
	l := token.NewList()
	l.PushStream(p.Receiver.Print(f))
	l.PushKeyword(token.ObjectOperator)
	l.PushStream(p.Name.Print(f))
	return l.String()
}
func (c ClassExpression) Print(f Format) string {
	l := token.NewList()
	l.PushStream(c.Receiver.Print(f))
	l.PushKeyword(token.ScopeResolutionOperator)
	l.PushStream(c.Expression.Print(f))
	return l
}
func (m Method) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(m.Visibility.Token())
	l.PushStream(m.FunctionStmt.Print(f))
	return l
}
func (m MethodCallExpression) Print(f Format) string {
	l := token.NewList()
	l.PushStream(m.Receiver.Print(f))
	l.PushKeyword(token.ObjectOperator)
	l.PushStream(m.FunctionCallExpression.Print(f))
	return l
}
func (i IfStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.If)
	l.PushKeyword(token.OpenParen)
	l.PushStream(i.Condition.Print(f))
	l.PushKeyword(token.CloseParen)
	l.PushStream(i.TrueBranch.Print(f))
	if i.FalseBranch != nil {
		l.PushKeyword(token.Else)
		l.PushStream(i.FalseBranch.Print(f))
	}
	return l
}
func (s SwitchStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Switch)
	l.PushKeyword(token.BlockBegin)
	for _, c := range s.Cases {
		l.PushStream(c.Print(f))
	}
	if s.DefaultCase != nil {
		l.PushKeyword(token.Default)
		l.PushKeyword(token.TernaryOperator2)
		l.PushStream(s.DefaultCase.Print(f))
	}
	l.PushKeyword(token.BlockEnd)
	return l
}
func (s SwitchCase) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Case)
	l.PushStream(s.Expression.Print(f))
	l.PushKeyword(token.TernaryOperator2)
	l.PushStream(s.Block.Print(f))
	return l
}
func (f ForStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.For)
	l.PushKeyword(token.OpenParen)
	for i, e := range f.Initialization {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(e.Print(f))
	}
	for i, e := range f.Termination {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(e.Print(f))
	}
	for i, e := range f.Iteration {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(e.Print(f))
	}
	l.PushKeyword(token.CloseParen)
	l.PushStream(f.LoopBlock.Print(f))
	return l
}
func (w WhileStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.While)
	l.PushKeyword(token.OpenParen)
	l.PushStream(w.Termination.Print(f))
	l.PushKeyword(token.CloseParen)
	l.PushStream(w.LoopBlock.Print(f))
	return l
}
func (w DoWhileStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Do)
	l.PushStream(w.LoopBlock.Print(f))
	l.PushKeyword(token.While)
	l.PushKeyword(token.OpenParen)
	l.PushStream(w.Termination.Print(f))
	l.PushKeyword(token.CloseParen)
	l.PushKeyword(token.StatementEnd)
	return l
}
func (t TryStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Try)
	l.PushStream(t.TryBlock.Print(f))
	for _, c := range t.CatchStmts {
		l.PushStream(c.Print(f))
	}
	if t.FinallyBlock != nil {
		l.PushKeyword(token.Finally)
		l.PushStream(t.FinallyBlock.Print(f))
	}
	return l
}
func (c CatchStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Catch)
	l.PushKeyword(token.OpenParen)
	l.Push(token.Item{Typ: token.Identifier, Val: c.CatchType})
	l.PushStream(c.CatchVar.Print(f))
	l.PushKeyword(token.CloseParen)
	l.PushStream(c.CatchBlock.Print(f))
	return l
}
func (l Literal) Print(f Format) string {
	switch l.Type {
	case String:
		return token.NewList(token.Item{Typ: token.StringLiteral, Val: l.Value})
	case Integer, Float:
		return token.NewList(token.Item{Typ: token.NumberLiteral, Val: l.Value})
	case Boolean:
		return token.NewList(token.Item{Typ: token.BooleanLiteral, Val: l.Value})
	case Null:
		return token.NewList(token.Item{Typ: token.Null, Val: "null"})
	}
	panic("invalid literal type")
}
func (f ForeachStmt) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Foreach)
	l.PushKeyword(token.OpenParen)
	l.PushStream(f.Source.Print(f))
	l.PushKeyword(token.AsOperator)
	if f.Key != nil {
		l.PushStream(f.Key.Print(f))
		l.PushKeyword(token.ArrayKeyOperator)
	}
	l.PushStream(f.Value.Print(f))
	l.PushKeyword(token.CloseParen)
	l.PushStream(f.LoopBlock.Print(f))
	return l
}
func (a ArrayExpression) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Array)
	l.PushKeyword(token.OpenParen)
	for i, pair := range a.Pairs {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(pair.Print(f))
	}
	return l
}
func (p ArrayPair) Print(f Format) string {
	l := token.NewList()
	if p.Key != nil {
		l.PushStream(p.Key.Print(f))
		l.PushKeyword(token.ArrayKeyOperator)
	}
	l.PushStream(p.Value.Print(f))
	return l
}
func (a ArrayLookupExpression) Print(f Format) string {
	l := token.NewList()
	l.PushStream(a.Array.Print(f))
	l.PushKeyword(token.ArrayLookupOperatorLeft)
	if a.Index != nil {
		l.PushStream(a.Index.Print(f))
	}
	l.PushKeyword(token.ArrayLookupOperatorRight)
	return l
}
func (a ArrayAppendExpression) Print(f Format) string {
	l := token.NewList()
	l.PushStream(a.Array.Print(f))
	l.PushKeyword(token.ArrayLookupOperatorLeft)
	l.PushKeyword(token.ArrayLookupOperatorRight)
	return l
}
func (s ShellCommand) Print(f Format) string {
	return token.NewList(token.Item{Typ: token.ShellCommand, Val: s.Command})
}
func (l ListStatement) Print(f Format) string {
	s := token.NewList()
	s.PushKeyword(token.List)
	for i, a := range l.Assignees {
		if i > 0 {
			s.PushKeyword(token.Comma)
		}
		s.PushStream(a.Print(f))
	}
	s.Push(token.Item{Typ: token.AssignmentOperator, Val: l.Operator})
	s.PushStream(l.Value.Print(f))
	return s
}

func (s StaticVariableDeclaration) Print(f Format) string {
	l := token.NewList()
	l.PushKeyword(token.Static)
	for i, d := range s.Declarations {
		if i > 0 {
			l.PushKeyword(token.Comma)
		}
		l.PushStream(d.Print(f))
	}
	l.PushKeyword(token.StatementEnd)
	return l
}
func (d DeclareBlock) Print(f Format) string {
	buf := bytes.NewBufferString("declare (")
	for i, decl := range d.Declarations {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(decl)
	}
	buf.WriteString(") {")
	buf.WriteString(d.Statements.Print(f))
	buf.WriteString("}")
	return buf.String()
}
