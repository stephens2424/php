package togo

import (
	"bytes"
	goast "go/ast"
	"go/token"
	"reflect"
	"strconv"

	phpast "github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/ast/printer"
)

func (t *Togo) ToGoStmt(php phpast.Statement) goast.Stmt {
	if v := reflect.ValueOf(php); v.Kind() == reflect.Ptr {
		php = v.Elem().Interface().(phpast.Statement)
	}

	switch n := php.(type) {
	// preliminary cases
	case phpast.UnaryExpression:
		if n.Operator == "--" || n.Operator == "++" {
			return &goast.IncDecStmt{
				X:   t.ToGoExpr(n.Operand),
				Tok: t.ToGoOperator(n.Operator),
			}
		}

	// standard cases
	case phpast.AnonymousFunction:
	case phpast.ArrayAppendExpression:
	case phpast.ArrayExpression:
	case phpast.ArrayLookupExpression:
	case phpast.AssignmentExpression:
		return &goast.AssignStmt{
			Lhs: []goast.Expr{t.ToGoExpr(n.Assignee)},
			Rhs: []goast.Expr{t.ToGoExpr(n.Value)},
			Tok: t.ToGoOperator(n.Operator),
		}
	case phpast.Block:
	case phpast.BreakStmt:
	case phpast.Class:
	case phpast.ClassExpression:
	case phpast.Constant:
	case phpast.ConstantExpression:
	case phpast.ContinueStmt:
	case phpast.DoWhileStmt:
	case phpast.EchoStmt:
		for _, e := range n.Expressions {
			return &goast.ExprStmt{t.CtxFuncCall("Echo.Write", []goast.Expr{t.ToGoExpr(e)})}
		}
	case phpast.EmptyStatement:
	case phpast.ExitStmt:
	case phpast.ExpressionStmt:
		switch expr := n.Expression.(type) {
		case phpast.AssignmentExpression:
			return t.ToGoStmt(expr)
		case *phpast.ShellCommand:
			return t.ToGoStmt(expr)
		case phpast.ShellCommand:
			return t.ToGoStmt(expr)
		}
		return &goast.ExprStmt{t.ToGoExpr(n.Expression)}
	case phpast.ForStmt:
		f := &goast.ForStmt{}
		if len(n.Initialization) == 1 {
			f.Init = t.ToGoStmt(n.Initialization[0])
		}

		// TODO Make sure all the termination expressions are *executed*, even though only the last one
		// is used to determine loop termination.
		if len(n.Termination) > 0 {
			f.Cond = t.ToGoExpr(n.Termination[len(n.Termination)-1])
		}
		f.Body = t.ToGoBlock(n.LoopBlock)

		// TODO Make sure all the iteration statements are *executed*
		if len(n.Iteration) > 0 {
			f.Post = t.ToGoStmt(n.Iteration[0])
		}
		return f
	case phpast.ForeachStmt:
		r := &goast.RangeStmt{}
		r.Key = t.ToGoExpr(n.Key)
		r.Value = t.ToGoExpr(n.Value)
		r.X = t.ToGoExpr(n.Source)
		r.Body = t.ToGoBlock(n.LoopBlock)
	case phpast.FunctionCallExpression:
	case phpast.FunctionCallStmt:
	case phpast.FunctionStmt:
	case phpast.GlobalDeclaration:
	case phpast.Identifier:
	case phpast.IfStmt:
		return t.TranslateIf(n)
	case phpast.Include:
	case phpast.IncludeStmt:
	case phpast.Interface:
	case phpast.ListStatement:
	case phpast.Method:
	case phpast.MethodCallExpression:
	case phpast.NewExpression:
	case phpast.PropertyExpression:
	case phpast.ReturnStmt:
	case phpast.ShellCommand:
		return &goast.ExprStmt{t.CtxFuncCall("Shell", []goast.Expr{&goast.BasicLit{Kind: token.STRING, Value: n.Command}})}
	case phpast.Statement:
	case phpast.StaticVariableDeclaration:
	case phpast.SwitchStmt:
	case phpast.ThrowStmt:
	case phpast.TryStmt:

	case phpast.Variable:
	case phpast.WhileStmt:
		f := &goast.ForStmt{}
		f.Cond = t.ToGoExpr(n.Termination)
		f.Body = t.ToGoBlock(n.LoopBlock)
		return f

		// broadest
	case phpast.Expression:
		return &goast.ExprStmt{t.ToGoExpr(n)}
	case phpast.Node:
	}

	return PHPEvalStmt(php)
}

func PHPEvalStmt(p phpast.Node) goast.Stmt {
	return &goast.ExprStmt{PHPEval(p)}
}

func PHPEval(p phpast.Node) goast.Expr {
	buf := &bytes.Buffer{}
	pr := printer.NewPrinter(buf)
	pr.PrintNode(p)
	return &goast.CallExpr{
		Fun: goast.NewIdent("PHPEval"),
		Args: []goast.Expr{
			&goast.BasicLit{Kind: token.STRING, Value: strconv.Quote(buf.String())},
		},
	}
}

func (t *Togo) ToGoExpr(p phpast.Expression) goast.Expr {
	if v := reflect.ValueOf(p); v.Kind() == reflect.Ptr {
		p = v.Elem().Interface().(phpast.Expression)
	}

	switch n := p.(type) {
	case phpast.AnonymousFunction:
	case phpast.ArrayAppendExpression:
	case phpast.ArrayExpression:
	case phpast.ArrayLookupExpression:
	case phpast.BinaryExpression:
		return &goast.BinaryExpr{
			X:  t.ToGoExpr(n.Antecedent),
			Y:  t.ToGoExpr(n.Subsequent),
			Op: t.ToGoOperator(n.Operator),
		}
	case phpast.UnaryExpression:
		return &goast.UnaryExpr{
			X:  t.ToGoExpr(n.Operand),
			Op: t.ToGoOperator(n.Operator),
		}
	case phpast.ClassExpression:
	case phpast.Constant:
	case phpast.ConstantExpression:
	case phpast.FunctionCallExpression:
	case phpast.Identifier:
		return goast.NewIdent(n.Value)
	case phpast.Include:
	case phpast.IncludeStmt:
	case phpast.Literal:
		switch n.Type {
		case phpast.String:
			return &goast.BasicLit{Kind: token.STRING, Value: n.Value}
		}
		return &goast.BasicLit{Value: n.Value}
	case phpast.MethodCallExpression:
	case phpast.NewExpression:
	case phpast.PropertyExpression:
		return t.ResolveDynamicProperty(t.ToGoExpr(n.Receiver), n.Name)
	case phpast.ShellCommand:
		return t.CtxFuncCall("Shell", []goast.Expr{&goast.BasicLit{Kind: token.STRING, Value: n.Command}})
	case phpast.Variable:
		return t.ToGoExpr(n.Name)
	}

	return PHPEval(p)
}

func (t *Togo) ToGoBlock(p phpast.Statement) *goast.BlockStmt {
	g := &goast.BlockStmt{}
	switch p := p.(type) {
	case *phpast.Block:
		g.List = []goast.Stmt{}
		for _, stmt := range p.Statements {
			g.List = append(g.List, t.ToGoStmt(stmt))
		}
	default:
		g.List = []goast.Stmt{t.ToGoStmt(p)}
	}
	return g
}

func (t *Togo) TranslateIf(p phpast.IfStmt) *goast.IfStmt {
	g := &goast.IfStmt{
		Cond: t.ToGoExpr(p.Branches[0].Condition),
		Body: t.ToGoBlock(p.Branches[0].Block),
	}

	if len(p.Branches) > 1 {
		g.Else = t.TranslateIf(phpast.IfStmt{
			Branches:  append([]phpast.IfBranch{}, p.Branches[1:]...),
			ElseBlock: p.ElseBlock,
		})
	}

	return g
}
