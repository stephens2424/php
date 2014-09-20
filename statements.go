package php

import (
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/token"
)

func (p *Parser) parseStmt() ast.Statement {
	switch p.current.typ {
	case token.BlockBegin:
		p.backup()
		return p.parseBlock()
	case token.Global:
		p.next()
		g := &ast.GlobalDeclaration{
			Identifiers: make([]*ast.Variable, 0, 1),
		}
		for p.current.typ == token.VariableOperator {
			variable, ok := p.parseVariable().(*ast.Variable)
			if !ok {
				p.errorf("global declarations must be of standard variables")
				break
			}
			g.Identifiers = append(g.Identifiers, variable)
			if p.peek().typ != token.Comma {
				break
			}
			p.expect(token.Comma)
			p.next()
		}
		p.expectStmtEnd()
		return g
	case token.Namespace:
		p.expect(token.Identifier)
		p.expectStmtEnd()
		// We are ignoring this for now
		return nil
	case token.Use:
		p.expect(token.Identifier)
		if p.peek().typ == token.AsOperator {
			p.expect(token.AsOperator)
			p.expect(token.Identifier)
		}
		p.expectStmtEnd()
		// We are ignoring this for now
		return nil
	case token.Static:
		if p.peek().typ == token.ScopeResolutionOperator {
			expr := p.parseExpression()
			p.expectStmtEnd()
			return expr
		}
		s := &ast.StaticVariableDeclaration{Declarations: make([]ast.Expression, 0)}
		for {
			p.expect(token.VariableOperator)
			p.expect(token.Identifier)
			v := ast.NewVariable(p.current.val)
			if p.peek().typ == token.AssignmentOperator {
				p.expect(token.AssignmentOperator)
				op := p.current.val
				p.expect(token.Null, token.StringLiteral, token.BooleanLiteral, token.NumberLiteral, token.Array)
				switch p.current.typ {
				case token.Array:
					s.Declarations = append(s.Declarations, &ast.AssignmentExpression{Assignee: v, Value: p.parseArrayDeclaration(), Operator: op})
				default:
					s.Declarations = append(s.Declarations, &ast.AssignmentExpression{Assignee: v, Value: p.parseLiteral(), Operator: op})
				}
			}
			s.Declarations = append(s.Declarations, v)
			if p.peek().typ != token.Comma {
				break
			}
			p.next()
		}
		p.expectStmtEnd()
		return s
	case token.VariableOperator, token.UnaryOperator:
		expr := ast.ExpressionStmt{p.parseExpression()}
		p.expectStmtEnd()
		return expr
	case token.Print:
		requireParen := false
		if p.peek().typ == token.OpenParen {
			p.expect(token.OpenParen)
			requireParen = true
		}
		stmt := ast.Echo(p.parseNextExpression())
		if requireParen {
			p.expect(token.CloseParen)
		}
		p.expectStmtEnd()
		return stmt
	case token.Function:
		return p.parseFunctionStmt()
	case token.PHPEnd:
		if p.peek().typ == token.EOF {
			return nil
		}
		var expr ast.Statement
		if p.accept(token.HTML) {
			expr = ast.Echo(&ast.Literal{Type: ast.String, Value: p.current.val})
		}
		p.next()
		if p.current.typ != token.EOF {
			p.expectCurrent(token.PHPBegin)
		}
		return expr
	case token.Echo:
		exprs := []ast.Expression{
			p.parseNextExpression(),
		}
		for p.peek().typ == token.Comma {
			p.expect(token.Comma)
			exprs = append(exprs, p.parseNextExpression())
		}
		p.expectStmtEnd()
		return ast.Echo(exprs...)
	case token.If:
		return p.parseIf()
	case token.While:
		return p.parseWhile()
	case token.Do:
		return p.parseDo()
	case token.For:
		return p.parseFor()
	case token.Foreach:
		return p.parseForeach()
	case token.Switch:
		return p.parseSwitch()
	case token.Abstract, token.Final, token.Class:
		return p.parseClass()
	case token.Interface:
		return p.parseInterface()
	case token.Return:
		p.next()
		stmt := ast.ReturnStmt{}
		if p.current.typ != token.StatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case token.Break:
		p.next()
		stmt := ast.BreakStmt{}
		if p.current.typ != token.StatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case token.Continue:
		p.next()
		stmt := ast.ContinueStmt{}
		if p.current.typ != token.StatementEnd {
			stmt.Expression = p.parseExpression()
			p.expectStmtEnd()
		}
		return stmt
	case token.Throw:
		stmt := ast.ThrowStmt{Expression: p.parseNextExpression()}
		p.expectStmtEnd()
		return stmt
	case token.Exit:
		stmt := ast.ExitStmt{}
		if p.peek().typ == token.OpenParen {
			p.expect(token.OpenParen)
			if p.peek().typ != token.CloseParen {
				stmt.Expression = p.parseNextExpression()
			}
			p.expect(token.CloseParen)
		}
		p.expectStmtEnd()
		return stmt
	case token.Try:
		stmt := &ast.TryStmt{}
		stmt.TryBlock = p.parseBlock()
		for p.expect(token.Catch); p.current.typ == token.Catch; p.next() {
			caught := &ast.CatchStmt{}
			p.expect(token.OpenParen)
			p.expect(token.Identifier)
			caught.CatchType = p.current.val
			p.expect(token.VariableOperator)
			p.expect(token.Identifier)
			caught.CatchVar = ast.NewVariable(p.current.val)
			p.expect(token.CloseParen)
			caught.CatchBlock = p.parseBlock()
			stmt.CatchStmts = append(stmt.CatchStmts, caught)
		}
		p.backup()
		return stmt
	case token.IgnoreErrorOperator:
		// Ignore this operator
		p.next()
		return p.parseStmt()
	case token.StatementEnd:
		// this is an empty statement
		return &ast.EmptyStatement{}
	case token.Declare:
		return p.parseDeclareBlock()
	default:
		expr := p.parseExpression()
		if expr != nil {
			p.expectStmtEnd()
			return ast.ExpressionStmt{expr}
		}
		p.errorf("Found %s, statement or expression", p.current)
		return nil
	}
}

func (p *Parser) expectStmtEnd() {
	if p.peek().typ != token.PHPEnd {
		p.expect(token.StatementEnd)
	}
}
