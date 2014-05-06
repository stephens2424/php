package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

func (p *Parser) parseFunctionStmt() *ast.FunctionStmt {
	stmt := &ast.FunctionStmt{}
	stmt.FunctionDefinition = p.parseFunctionDefinition()
	stmt.Body = p.parseBlock()
	return stmt
}

func (p *Parser) parseFunctionDefinition() *ast.FunctionDefinition {
	def := &ast.FunctionDefinition{}
	if p.peek().typ == token.AmpersandOperator {
		// This is a function returning a reference ... ignore this for now
		p.next()
	}
	if !p.accept(token.Identifier) {
		p.next()
		if !isKeyword(p.current.typ, p.current.val) {
			p.errorf("bad function name", p.current.val)
		}
	}
	def.Name = p.current.val
	def.Arguments = make([]ast.FunctionArgument, 0)
	p.expect(token.OpenParen)
	if p.peek().typ == token.CloseParen {
		p.expect(token.CloseParen)
		return def
	}
	def.Arguments = append(def.Arguments, p.parseFunctionArgument())
	for {
		switch p.peek().typ {
		case token.Comma:
			p.expect(token.Comma)
			def.Arguments = append(def.Arguments, p.parseFunctionArgument())
		case token.CloseParen:
			p.expect(token.CloseParen)
			return def
		default:
			p.errorf("unexpected argument separator:", p.current)
			return def
		}
	}
}

func (p *Parser) parseFunctionArgument() ast.FunctionArgument {
	arg := ast.FunctionArgument{}
	switch p.peek().typ {
	case token.Identifier, token.Array:
		p.next()
		arg.TypeHint = p.current.val
	}
	if p.peek().typ == token.AmpersandOperator {
		p.next()
	}
	p.expect(token.VariableOperator)
	p.next()
	arg.Variable = ast.NewVariable(p.current.val)
	if p.peek().typ == token.AssignmentOperator {
		p.expect(token.AssignmentOperator)
		p.next()
		arg.Default = p.parseExpression()
	}
	return arg
}

func (p *Parser) parseFunctionCall(callable ast.Expression) *ast.FunctionCallExpression {
	expr := &ast.FunctionCallExpression{}
	expr.FunctionName = callable
	return p.parseFunctionArguments(expr)
}

func (p *Parser) parseFunctionArguments(expr *ast.FunctionCallExpression) *ast.FunctionCallExpression {
	expr.Arguments = make([]ast.Expression, 0)
	p.expect(token.OpenParen)
	if p.peek().typ == token.CloseParen {
		p.expect(token.CloseParen)
		return expr
	}
	expr.Arguments = append(expr.Arguments, p.parseNextExpression())
	for p.peek().typ != token.CloseParen {
		p.expect(token.Comma)
		arg := p.parseNextExpression()
		if arg == nil {
			break
		}
		expr.Arguments = append(expr.Arguments, arg)
	}
	p.expect(token.CloseParen)
	return expr

}

func (p *Parser) parseAnonymousFunction() ast.Expression {
	f := &ast.AnonymousFunction{}
	f.Arguments = make([]ast.FunctionArgument, 0)
	f.ClosureVariables = make([]ast.FunctionArgument, 0)
	p.expect(token.OpenParen)
	if p.peek().typ != token.CloseParen {
		f.Arguments = append(f.Arguments, p.parseFunctionArgument())
	}

Loop:
	for {
		switch p.peek().typ {
		case token.Comma:
			p.expect(token.Comma)
			f.Arguments = append(f.Arguments, p.parseFunctionArgument())
		case token.CloseParen:
			break Loop
		default:
			p.errorf("unexpected argument separator:", p.current)
			return f
		}
	}
	p.expect(token.CloseParen)

	// Closure variables
	if p.peek().typ == token.Use {
		p.expect(token.Use)
		p.expect(token.OpenParen)
		f.ClosureVariables = append(f.ClosureVariables, p.parseFunctionArgument())
	ClosureLoop:
		for {
			switch p.peek().typ {
			case token.Comma:
				p.expect(token.Comma)
				f.ClosureVariables = append(f.ClosureVariables, p.parseFunctionArgument())
			case token.CloseParen:
				break ClosureLoop
			default:
				p.errorf("unexpected argument separator:", p.current)
				return f
			}
		}
		p.expect(token.CloseParen)
	}

	f.Body = p.parseBlock()
	return f
}
