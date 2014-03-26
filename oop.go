package php

import "stephensearles.com/php/ast"

func (p *parser) parseInstantiation() ast.Expression {
	p.expectCurrent(itemNewOperator)
	expr := &ast.NewExpression{}
	expr.Class = p.parseNextExpression()

	if p.peek().typ == itemOpenParen {
		p.expect(itemOpenParen)
		if p.peek().typ != itemCloseParen {
			expr.Arguments = append(expr.Arguments, p.parseNextExpression())
			for p.peek().typ == itemComma {
				p.expect(itemComma)
				expr.Arguments = append(expr.Arguments, p.parseNextExpression())
			}
		}
		p.expect(itemCloseParen)
	}
	return expr
}

func (p *parser) parseClass() ast.Class {
	if p.current.typ == itemAbstract {
		p.expect(itemClass)
	}
	p.expect(itemIdentifier)
	name := p.current.val
	if p.peek().typ == itemExtends {
		p.expect(itemExtends)
		p.expect(itemIdentifier)
	}
	if p.peek().typ == itemImplements {
		p.expect(itemImplements)
		p.expect(itemIdentifier)
		for p.peek().typ == itemComma {
			p.expect(itemComma)
			p.expect(itemIdentifier)
		}
	}
	p.expect(itemBlockBegin)
	return p.parseClassFields(ast.Class{Name: name})
}

func (p *parser) parseObjectLookup(r ast.Expression) (expr ast.Expression) {
	p.expectCurrent(itemObjectOperator)
	prop := &ast.PropertyExpression{
		Receiver: r,
	}
	switch p.next(); p.current.typ {
	case itemBlockBegin:
		prop.Name = p.parseNextExpression()
		p.expect(itemBlockEnd)
	case itemVariableOperator:
		prop.Name = p.parseExpression()
	case itemIdentifier:
		prop.Name = ast.Identifier{Value: p.current.val}
	}
	expr = prop
	switch pk := p.peek(); pk.typ {
	case itemOpenParen:
		expr = &ast.MethodCallExpression{
			Receiver:               r,
			FunctionCallExpression: p.parseFunctionCall(prop.Name),
		}
	}
	expr = p.parseOperation(p.parenLevel, expr)
	return
}

func (p *parser) parseVisibility() (vis ast.Visibility, found bool) {
	switch p.peek().typ {
	case itemPrivate:
		vis = ast.Private
	case itemPublic:
		vis = ast.Public
	case itemProtected:
		vis = ast.Protected
	default:
		return ast.Public, false
	}
	p.next()
	return vis, true
}

func (p *parser) parseAbstract() bool {
	if p.peek().typ == itemAbstract {
		p.next()
		return true
	}
	return false
}

func (p *parser) parseClassFields(c ast.Class) ast.Class {
	c.Methods = make([]ast.Method, 0)
	c.Properties = make([]ast.Property, 0)
	for p.peek().typ != itemBlockEnd {
		vis, _, _, abstract := p.parseClassMemberSettings()
		p.next()
		switch p.current.typ {
		case itemFunction:
			if abstract {
				f := p.parseFunctionDefinition()
				m := ast.Method{
					Visibility:   vis,
					FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
				}
				c.Methods = append(c.Methods, m)
				p.expect(itemStatementEnd)
			} else {
				c.Methods = append(c.Methods, ast.Method{
					Visibility:   vis,
					FunctionStmt: p.parseFunctionStmt(),
				})
			}
		case itemVar:
			p.next()
			fallthrough
		case itemVariableOperator:
			p.next()
			prop := ast.Property{
				Visibility: vis,
				Name:       "$" + p.current.val,
			}
			if p.peek().typ == itemAssignmentOperator {
				p.expect(itemAssignmentOperator)
				prop.Initialization = p.parseNextExpression()
			}
			c.Properties = append(c.Properties, prop)
			p.expect(itemStatementEnd)
		case itemConst:
			constant := ast.Constant{}
			p.expect(itemIdentifier)
			constant.Variable = ast.NewVariable(p.current.val)
			if p.peek().typ == itemAssignmentOperator {
				p.expect(itemAssignmentOperator)
				constant.Value = p.parseNextExpression()
			}
			c.Constants = append(c.Constants, constant)
			p.expect(itemStatementEnd)
		default:
			p.errorf("unexpected class member %v", p.current)
		}
	}
	p.expect(itemBlockEnd)
	return c
}

func (p *parser) parseInterface() *ast.Interface {
	i := &ast.Interface{
		Inherits: make([]string, 0),
	}
	p.expect(itemIdentifier)
	i.Name = p.current.val
	if p.peek().typ == itemExtends {
		p.expect(itemExtends)
		for {
			p.expect(itemIdentifier)
			i.Inherits = append(i.Inherits, p.current.val)
			if p.peek().typ != itemComma {
				break
			}
			p.expect(itemComma)
		}
	}
	p.expect(itemBlockBegin)
	for p.peek().typ != itemBlockEnd {
		vis, _ := p.parseVisibility()
		if p.peek().typ == itemStatic {
			p.next()
		}
		p.next()
		switch p.current.typ {
		case itemFunction:
			f := p.parseFunctionDefinition()
			m := ast.Method{
				Visibility:   vis,
				FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
			}
			i.Methods = append(i.Methods, m)
			p.expect(itemStatementEnd)
		default:
			p.errorf("unexpected interface member %v", p.current)
		}
	}
	p.expect(itemBlockEnd)
	return i
}

func (p *parser) parseClassMemberSettings() (vis ast.Visibility, static, final, abstract bool) {
	var foundVis bool
	vis = ast.Public
	for {
		switch p.peek().typ {
		case itemAbstract:
			if abstract {
				p.errorf("found multiple abstract declarations")
			}
			abstract = true
			p.next()
		case itemPrivate, itemPublic, itemProtected:
			if foundVis {
				p.errorf("found multiple visibility declarations")
			}
			vis, foundVis = p.parseVisibility()
		case itemFinal:
			if final {
				p.errorf("found multiple final declarations")
			}
			final = true
			p.next()
		case itemStatic:
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
