package togo

import (
	goast "go/ast"
	"go/token"
	"reflect"
	"strconv"

	phpast "github.com/stephens2424/php/ast"
)

func ToGo(php phpast.Node) goast.Node {
	if v := reflect.ValueOf(php); v.Kind() == reflect.Ptr {
		php = v.Elem().Interface().(phpast.Node)
	}

	switch n := php.(type) {
	case phpast.AnonymousFunction:
	case phpast.ArrayAppendExpression:
	case phpast.ArrayExpression:
	case phpast.ArrayLookupExpression:
	case phpast.ArrayPair:
	case phpast.Assignable:
	case phpast.AssignmentExpression:
	case phpast.Block:
	case phpast.BreakStmt:
	case phpast.BinaryExpression:
		return &goast.BinaryExpr{
			X:  ToGoExpr(n.Antecedent),
			Y:  ToGoExpr(n.Subsequent),
			Op: ToGoOperator(n.Operator),
		}
	case phpast.CatchStmt:
	case phpast.Class:
	case phpast.ClassExpression:
	case phpast.Constant:
	case phpast.ConstantExpression:
	case phpast.ContinueStmt:
	case phpast.DoWhileStmt:
	case phpast.EchoStmt:
	case phpast.EmptyStatement:
	case phpast.ExitStmt:
	case phpast.ExpressionStmt:
	case phpast.ForStmt:
		f := &goast.ForStmt{}
		if len(n.Initialization) == 1 {
			f.Init = ToGoStmt(n.Initialization[0])
		}

		// TODO Make sure all the termination expressions are *executed*, even though only the last one
		// is used to determine loop termination.
		if len(n.Termination) > 0 {
			f.Cond = ToGoExpr(n.Termination[len(n.Termination)-1])
		}
		f.Body = ToGoBlock(n.LoopBlock)

		// TODO Make sure all the iteration statements are *executed*
		if len(n.Iteration) > 0 {
			f.Post = ToGoStmt(n.Iteration[0])
		}
		return f
	case phpast.ForeachStmt:
		r := &goast.RangeStmt{}
		r.Key = ToGoExpr(n.Key)
		r.Value = ToGoExpr(n.Value)
		r.X = ToGoExpr(n.Source)
		r.Body = ToGoBlock(n.LoopBlock)
	case phpast.FunctionArgument:
	case phpast.FunctionCallExpression:
	case phpast.FunctionCallStmt:
	case phpast.FunctionDefinition:
	case phpast.FunctionStmt:
	case phpast.GlobalDeclaration:
	case phpast.Identifier:
	case phpast.IfStmt:
		return TranslateIf(n)
	case phpast.Include:
	case phpast.IncludeStmt:
	case phpast.Interface:
	case phpast.ListStatement:
	case phpast.Literal:
		switch n.Type {
		case phpast.String:
			return &goast.BasicLit{Kind: token.STRING, Value: n.Value}
		default:
			return &goast.BasicLit{Value: n.Value}
		}
	case phpast.Method:
	case phpast.MethodCallExpression:
	case phpast.NewExpression:
	case phpast.Node:
	case phpast.Property:
	case phpast.PropertyExpression:
	case phpast.ReturnStmt:
	case phpast.ShellCommand:
	case phpast.Statement:
	case phpast.StaticVariableDeclaration:
	case phpast.SwitchCase:
	case phpast.SwitchStmt:
	case phpast.ThrowStmt:
	case phpast.TryStmt:
	case phpast.Variable:
	case phpast.WhileStmt:
		f := &goast.ForStmt{}
		f.Cond = ToGoExpr(n.Termination)
		f.Body = ToGoBlock(n.LoopBlock)
		return f
	}

	return PHPEval(php)
}

func PHPEval(p phpast.Node) goast.Node {
	return &goast.CallExpr{
		Fun: goast.NewIdent("PHPEval"),
		Args: []goast.Expr{
			&goast.BasicLit{Kind: token.STRING, Value: strconv.Quote(p.String())},
		},
	}
}

func ToGoExpr(p phpast.Expression) goast.Expr {
	if e := ToGo(p); e != nil {
		return e.(goast.Expr)
	}
	return nil
}

func ToGoBlock(p phpast.Statement) *goast.BlockStmt {
	g := &goast.BlockStmt{}
	switch p := p.(type) {
	case *phpast.Block:
		g.List = []goast.Stmt{}
		for _, stmt := range p.Statements {
			g.List = append(g.List, ToGoStmt(stmt))
		}
	default:
		g.List = []goast.Stmt{ToGoStmt(p)}
	}
	return g
}

func ToGoStmt(p phpast.Statement) goast.Stmt {
	return ToGo(p).(goast.Stmt)
}

func TranslateIf(p phpast.IfStmt) *goast.IfStmt {
	g := &goast.IfStmt{
		Cond: ToGoExpr(p.Branches[0].Condition),
		Body: ToGoBlock(p.Branches[0].Block),
	}

	if len(p.Branches) > 1 {
		g.Else = TranslateIf(phpast.IfStmt{
			Branches:  append([]phpast.IfBranch{}, p.Branches[1:]...),
			ElseBlock: p.ElseBlock,
		})
	}

	return g
}
