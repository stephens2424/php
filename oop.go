package php

import (
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/lexer"
	"github.com/stephens2424/php/token"
)

func (p *Parser) parseInstantiation() ast.Expression {
	p.expectCurrent(token.NewOperator)
	p.next()

	p.instantiation = true
	expr := &ast.NewExpression{}
	expr.Class = p.parseOperand()
	p.instantiation = false

	if p.peek().Typ == token.OpenParen {
		p.expect(token.OpenParen)
		if p.peek().Typ != token.CloseParen {
			expr.Arguments = append(expr.Arguments, p.parseNextExpression())
			for p.peek().Typ == token.Comma {
				p.expect(token.Comma)
				expr.Arguments = append(expr.Arguments, p.parseNextExpression())
			}
		}
		p.expect(token.CloseParen)
	}
	return expr
}

func (p *Parser) parseClass() *ast.Class {
	if p.current.Typ == token.Abstract {
		p.expect(token.Class)
	}
	if p.current.Typ == token.Final {
		p.expect(token.Class)
	}
	switch p.next(); {
	case p.current.Typ == token.Identifier:
	case lexer.IsKeyword(p.current.Typ, p.current.Val):
	default:
		p.errorf("unexpected variable operand %s", p.current)
	}

	name := p.current.Val
	if p.peek().Typ == token.Extends {
		p.expect(token.Extends)
		p.expect(token.Identifier)
	}
	if p.peek().Typ == token.Implements {
		p.expect(token.Implements)
		p.expect(token.Identifier)
		for p.peek().Typ == token.Comma {
			p.expect(token.Comma)
			p.expect(token.Identifier)
		}
	}
	p.expect(token.BlockBegin)
	return p.parseClassFields(&ast.Class{Name: name})
}

func (p *Parser) parseObjectLookup(r ast.Expression) (expr ast.Expression) {
	p.expectCurrent(token.ObjectOperator)
	prop := &ast.PropertyExpression{
		Receiver: r,
	}
	switch p.next(); p.current.Typ {
	case token.BlockBegin:
		prop.Name = p.parseNextExpression()
		p.expect(token.BlockEnd)
	case token.VariableOperator:
		prop.Name = p.parseExpression()
	case token.Identifier:
		prop.Name = ast.Identifier{Value: p.current.Val}
	}
	expr = prop
	switch pk := p.peek(); pk.Typ {
	case token.OpenParen:
		expr = &ast.MethodCallExpression{
			Receiver:               r,
			FunctionCallExpression: p.parseFunctionCall(prop.Name),
		}
	}
	expr = p.parseOperation(p.parenLevel, expr)
	return
}

func (p *Parser) parseVisibility() (vis ast.Visibility, found bool) {
	switch p.peek().Typ {
	case token.Private:
		vis = ast.Private
	case token.Public:
		vis = ast.Public
	case token.Protected:
		vis = ast.Protected
	default:
		return ast.Public, false
	}
	p.next()
	return vis, true
}

func (p *Parser) parseAbstract() bool {
	if p.peek().Typ == token.Abstract {
		p.next()
		return true
	}
	return false
}

func (p *Parser) parseClassFields(c *ast.Class) *ast.Class {
	// Starting on BlockBegin
	c.Methods = make([]ast.Method, 0)
	c.Properties = make([]ast.Property, 0)
	for p.peek().Typ != token.BlockEnd {
		vis, _, _, abstract := p.parseClassMemberSettings()
		p.next()
		switch p.current.Typ {
		case token.Function:
			if abstract {
				f := p.parseFunctionDefinition()
				m := ast.Method{
					Visibility:   vis,
					FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
				}
				c.Methods = append(c.Methods, m)
				p.expect(token.StatementEnd)
			} else {
				c.Methods = append(c.Methods, ast.Method{
					Visibility:   vis,
					FunctionStmt: p.parseFunctionStmt(),
				})
			}
		case token.Var:
			p.expect(token.VariableOperator)
			fallthrough
		case token.VariableOperator:
			for {
				p.expect(token.Identifier)
				prop := ast.Property{
					Visibility: vis,
					Name:       "$" + p.current.Val,
				}
				if p.peek().Typ == token.AssignmentOperator {
					p.expect(token.AssignmentOperator)
					prop.Initialization = p.parseNextExpression()
				}
				c.Properties = append(c.Properties, prop)
				if p.accept(token.StatementEnd) {
					break
				}
				p.expect(token.Comma)
				p.expect(token.VariableOperator)
			}
		case token.Const:
			constant := ast.Constant{}
			p.expect(token.Identifier)
			constant.Variable = ast.NewVariable(p.current.Val)
			if p.peek().Typ == token.AssignmentOperator {
				p.expect(token.AssignmentOperator)
				constant.Value = p.parseNextExpression()
			}
			c.Constants = append(c.Constants, constant)
			p.expect(token.StatementEnd)
		default:
			p.errorf("unexpected class member %v", p.current)
			return c
		}
	}
	p.expect(token.BlockEnd)
	return c
}

func (p *Parser) parseInterface() *ast.Interface {
	i := &ast.Interface{
		Inherits: make([]string, 0),
	}
	p.expect(token.Identifier)
	i.Name = p.current.Val
	if p.peek().Typ == token.Extends {
		p.expect(token.Extends)
		for {
			p.expect(token.Identifier)
			i.Inherits = append(i.Inherits, p.current.Val)
			if p.peek().Typ != token.Comma {
				break
			}
			p.expect(token.Comma)
		}
	}
	p.expect(token.BlockBegin)
	for p.peek().Typ != token.BlockEnd {
		vis, _ := p.parseVisibility()
		if p.peek().Typ == token.Static {
			p.next()
		}
		p.next()
		switch p.current.Typ {
		case token.Function:
			f := p.parseFunctionDefinition()
			m := ast.Method{
				Visibility:   vis,
				FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
			}
			i.Methods = append(i.Methods, m)
			p.expect(token.StatementEnd)
		case token.Const:
			constant := ast.Constant{}
			p.expect(token.Identifier)
			constant.Variable = ast.NewVariable(p.current.Val)
			if p.peek().Typ == token.AssignmentOperator {
				p.expect(token.AssignmentOperator)
				constant.Value = p.parseNextExpression()
			}
			i.Constants = append(i.Constants, constant)
			p.expect(token.StatementEnd)
		default:
			p.errorf("unexpected interface member %v", p.current)
		}
	}
	p.expect(token.BlockEnd)
	return i
}

func (p *Parser) parseClassMemberSettings() (vis ast.Visibility, static, final, abstract bool) {
	var foundVis bool
	vis = ast.Public
	for {
		switch p.peek().Typ {
		case token.Abstract:
			if abstract {
				p.errorf("found multiple abstract declarations")
			}
			abstract = true
			p.next()
		case token.Private, token.Public, token.Protected:
			if foundVis {
				p.errorf("found multiple visibility declarations")
			}
			vis, foundVis = p.parseVisibility()
		case token.Final:
			if final {
				p.errorf("found multiple final declarations")
			}
			final = true
			p.next()
		case token.Static:
			if static {
				p.errorf("found multiple static declarations")
			}
			static = true
			p.next()
		default:
			return
		}
	}
	return
}
