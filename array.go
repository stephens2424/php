package php

import (
	"github.com/stephens2424/php/ast"
	"github.com/stephens2424/php/token"
)

func (p *Parser) parseArrayLookup(e ast.Expression) ast.Expression {
	p.expectCurrent(token.ArrayLookupOperatorLeft, token.BlockBegin)
	switch Typ := p.peek().Typ; Typ {
	case token.ArrayLookupOperatorRight, token.BlockBegin:
		p.expect(token.ArrayLookupOperatorRight, token.BlockEnd)
		return ast.ArrayAppendExpression{Array: e}
	}
	p.next()
	expr := &ast.ArrayLookupExpression{
		Array: e,
		Index: p.parseExpression(),
	}
	p.expect(token.ArrayLookupOperatorRight, token.BlockEnd)
	return expr
}

func (p *Parser) parseArrayDeclaration() ast.Expression {
	var endType token.Token
	var pairs []ast.ArrayPair
	p.expectCurrent(token.Array, token.ArrayLookupOperatorLeft)
	switch p.current.Typ {
	case token.Array:
		p.expect(token.OpenParen)
		endType = token.CloseParen
	case token.ArrayLookupOperatorLeft:
		endType = token.ArrayLookupOperatorRight
	}
ArrayLoop:
	for {
		var key, Val ast.Expression
		switch p.peek().Typ {
		case endType:
			break ArrayLoop
		default:
			Val = p.parseNextExpression()
		}
		switch p.peek().Typ {
		case token.Comma:
			p.expect(token.Comma)
		case endType:
			pairs = append(pairs, ast.ArrayPair{Key: key, Value: Val})
			break ArrayLoop
		case token.ArrayKeyOperator:
			p.expect(token.ArrayKeyOperator)
			key = Val
			Val = p.parseNextExpression()
			if p.peek().Typ == endType {
				pairs = append(pairs, ast.ArrayPair{Key: key, Value: Val})
				break ArrayLoop
			}
			p.expect(token.Comma)
		default:
			p.errorf("expected => or ,")
			return nil
		}
		pairs = append(pairs, ast.ArrayPair{Key: key, Value: Val})
	}
	p.expect(endType)
	return &ast.ArrayExpression{Pairs: pairs}
}

func (p *Parser) parseList() ast.Expression {
	l := &ast.ListStatement{
		Assignees: make([]ast.Assignable, 0),
	}
	p.expect(token.OpenParen)
	for {
		if p.accept(token.Comma) {
			continue
		}
		if p.peek().Typ == token.CloseParen {
			break
		}
		p.next()
		op, ok := p.parseOperand().(ast.Assignable)
		if ok {
			l.Assignees = append(l.Assignees, op)
		} else {
			p.errorf("%v list element is not assignable", op)
		}
		if p.peek().Typ != token.Comma {
			break
		}
		p.expect(token.Comma)
	}
	p.expect(token.CloseParen)
	p.expect(token.AssignmentOperator)
	l.Operator = p.current.Val
	l.Value = p.parseNextExpression()
	return l

}
