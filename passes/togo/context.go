package togo

import (
	goast "go/ast"

	phpast "github.com/stephens2424/php/ast"
)

type context struct {
	Scope phpast.Scope
}

func (t *Togo) ResolveDynamicVar(varName phpast.Expression) goast.Node {
	switch e := varName.(type) {
	case phpast.Identifier:
		return goast.NewIdent(e.Value)
	}

	return t.CtxFuncCall("GetDynamic", []goast.Expr{t.ToGoExpr(varName)})
}

func (t *Togo) ResolveDynamicProperty(rcvr goast.Expr, propName phpast.Expression) goast.Expr {
	switch e := propName.(type) {
	case phpast.Identifier:
		return &goast.SelectorExpr{
			X:   rcvr,
			Sel: goast.NewIdent(e.Value),
		}
	}

	return t.CtxFuncCall("GetDynamicProperty", []goast.Expr{rcvr, t.ToGoExpr(propName)})
}

func (t *Togo) CtxFuncCall(funcName string, args []goast.Expr) *goast.CallExpr {
	return &goast.CallExpr{
		Fun: &goast.SelectorExpr{
			X:   goast.NewIdent("phpctx"),
			Sel: goast.NewIdent(funcName),
		},
		Args: args,
	}
}
